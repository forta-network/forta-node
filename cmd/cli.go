package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"

	"github.com/forta-network/forta-node/config"
	"gopkg.in/yaml.v3"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyFortaDir         = "forta_dir"
	keyFortaPassphrase  = "forta_passphrase"
	keyFortaDevelopment = "forta_development"
	keyFortaExposeNats  = "forta_expose_nats"
)

var (
	cfg config.Config

	parsedArgs struct {
		Version uint64
		NoCheck bool
	}

	cmdForta = &cobra.Command{
		Use:   "forta",
		Short: "Forta node command line interface",
		Long: `Forta node scans blockchains by using agents to detect vulnerabilities and
publishes alerts about them`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage: true,
	}

	cmdFortaInit = &cobra.Command{
		Use:   "init",
		Short: "initialize a config file and a private key (doesn't overwrite)",
		RunE:  handleFortaInit,
	}

	cmdFortaRun = &cobra.Command{
		Use:   "run",
		Short: "launch the node",
		RunE:  withInitialized(withValidConfig(handleFortaRun)),
	}

	cmdFortaAccount = &cobra.Command{
		Use:   "account",
		Short: "account management",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmdFortaAccountAddress = &cobra.Command{
		Use:   "address",
		Short: "show the scanner address",
		RunE:  withInitialized(handleFortaAccountAddress),
	}

	cmdFortaAccountImport = &cobra.Command{
		Use:    "import",
		Short:  "import new scanner account (removes the old one)",
		RunE:   withInitialized(handleFortaAccountImport),
		Hidden: true,
	}

	cmdFortaImages = &cobra.Command{
		Use:   "images",
		Short: "list the Forta node container images",
		RunE:  handleFortaImages,
	}

	cmdFortaVersion = &cobra.Command{
		Use:   "version",
		Short: "show release info",
		RunE:  handleFortaVersion,
	}

	cmdFortaBatch = &cobra.Command{
		Use:   "batch",
		Short: "batch utils",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmdFortaStatus = &cobra.Command{
		Use:   "status",
		Short: "display statuses of node services",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			show, err := cmd.Flags().GetString("show")
			if err != nil {
				return err
			}
			noColor, err := cmd.Flags().GetBool("no-color")
			if err != nil {
				return err
			}
			return handleFortaStatus(cmd, format, show, noColor)
		},
	}

	cmdFortaStatusAll = &cobra.Command{
		Use:   "all",
		Short: "shorthand for `--show all --format oneline`",
		RunE: func(cmd *cobra.Command, args []string) error {
			noColor, err := cmd.Flags().GetBool("no-color")
			if err != nil {
				return err
			}
			return handleFortaStatus(cmd, StatusFormatOneline, StatusShowAll, noColor)
		},
	}

	cmdFortaAuthorize = &cobra.Command{
		Use:   "authorize",
		Short: "generate a signature for a specific action",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmdFortaAuthorizePool = &cobra.Command{
		Use:   "pool",
		Short: "generate a pool registration signature",
		RunE:  withInitialized(withValidConfig(handleFortaAuthorizePool)),
	}
)

// Execute executes the root command.
func Execute() error {
	return cmdForta.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	cmdForta.AddCommand(cmdFortaInit)
	cmdForta.AddCommand(cmdFortaRun)

	cmdForta.AddCommand(cmdFortaAccount)
	cmdFortaAccount.AddCommand(cmdFortaAccountAddress)
	cmdFortaAccount.AddCommand(cmdFortaAccountImport)

	cmdForta.AddCommand(cmdFortaImages)

	cmdForta.AddCommand(cmdFortaVersion)

	cmdForta.AddCommand(cmdFortaBatch)

	cmdForta.AddCommand(cmdFortaStatus)
	cmdFortaStatus.AddCommand(cmdFortaStatusAll)

	cmdForta.AddCommand(cmdFortaAuthorize)
	cmdFortaAuthorize.AddCommand(cmdFortaAuthorizePool)

	// Global (persistent) flags

	cmdForta.PersistentFlags().String("dir", "", "Forta dir (default is $HOME/.forta) (overrides $FORTA_DIR)")
	viper.BindPFlag(keyFortaDir, cmdForta.PersistentFlags().Lookup("dir"))

	cmdForta.PersistentFlags().Bool("development", false, "development mode (overrides $FORTA_DEVELOPMENT)")
	viper.BindPFlag(keyFortaDevelopment, cmdForta.PersistentFlags().Lookup("development"))

	cmdForta.PersistentFlags().String("passphrase", "", "passphrase to decrypt the private key (overrides $FORTA_PASSPHRASE)")
	viper.BindPFlag(keyFortaPassphrase, cmdForta.PersistentFlags().Lookup("passphrase"))

	cmdForta.PersistentFlags().Bool("expose-nats", false, "expose nats via public docker network")
	viper.BindPFlag(keyFortaExposeNats, cmdForta.PersistentFlags().Lookup("expose-nats"))

	// forta account import
	cmdFortaAccountImport.Flags().String("file", "", "path to a file that contains a private key hex")
	cmdFortaAccountImport.MarkFlagRequired("file")

	// forta run
	cmdFortaRun.Flags().BoolVar(&parsedArgs.NoCheck, "no-check", false, "disable scanner registry check and just run")

	// forta status
	cmdFortaStatus.Flags().String("format", StatusFormatPretty, "output formatting/encoding: pretty (default), oneline, json, csv")
	cmdFortaStatus.Flags().Bool("no-color", false, "disable colors")
	cmdFortaStatus.Flags().String("show", StatusShowSummary, "filter statuses to show: summary (default), important, all")

	// forta status all
	cmdFortaStatusAll.Flags().Bool("no-color", false, "disable colors")

	// forta authorize pool
	cmdFortaAuthorizePool.Flags().String("id", "", "scanner pool ID (integer)")
	cmdFortaAuthorizePool.MarkFlagRequired("id")
	cmdFortaAuthorizePool.Flags().Bool("polygonscan", false, "see the registerScannerNode() inputs to use in Polygonscan")
	cmdFortaAuthorizePool.Flags().BoolP("force", "f", false, "ignore warning(s)")
	cmdFortaAuthorizePool.Flags().Bool("clean", false, "output only the encoded registration info")
}

func initConfig() {
	viper.SetConfigType("yaml")

	viper.BindEnv(keyFortaDir)
	viper.BindEnv(keyFortaPassphrase)
	viper.BindEnv(keyFortaDevelopment)
	viper.BindEnv(keyFortaExposeNats)
	viper.AutomaticEnv()

	fortaDir := viper.GetString(keyFortaDir)
	if fortaDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			logrus.Panicf("failed to get home dir: %v", err)
		}
		fortaDir = path.Join(home, ".forta")
	}

	configPath := path.Join(fortaDir, config.DefaultConfigFileName)
	configBytes, _ := ioutil.ReadFile(configPath)
	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		yellowBold("Your config file is invalid! Please check the values and fix any formatting issues.\n")
		logrus.WithError(err).Fatal("failed to read config")
	}

	if err := defaults.Set(&cfg); err != nil {
		panic(err)
	}

	cfg.FortaDir = fortaDir
	cfg.KeyDirPath = path.Join(cfg.FortaDir, config.DefaultKeysDirName)
	cfg.Development = viper.GetBool(keyFortaDevelopment)
	cfg.Passphrase = viper.GetString(keyFortaPassphrase)

	viper.ReadConfig(bytes.NewBuffer(configBytes))
	config.InitLogLevel(cfg)
}

var configEnvVarRegexp = regexp.MustCompile(`\$[A-Z0-9_]+`)

// TODO: viper.Unmarshal is a mess. Use this again somehow with a custom hook.
func resolveEnvVarsHookFunc(from reflect.Kind, to reflect.Kind, data interface{}) (interface{}, error) {
	if from != reflect.String {
		return data, nil
	}

	// Resolve $ENV_VAR expressions to the values from the host OS environment.
	return configEnvVarRegexp.ReplaceAllStringFunc(data.(string), func(match string) string {
		return os.Getenv(match[1:])
	}), nil
}

func validateConfig() error {
	validate := validator.New()

	// Use the YAML names while validating the struct.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := validate.Struct(&cfg); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		fmt.Fprintln(os.Stderr, "The config file has invalid or missing fields:")
		for _, validationErr := range validationErrs {
			fmt.Fprintf(os.Stderr, "  - %s\n", validationErr.Namespace()[7:])
		}
		return errors.New("invalid config file")
	}

	return nil
}

func withValidConfig(handler func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := validateConfig(); err != nil {
			return err
		}
		return handler(cmd, args)
	}
}

func withInitialized(handler func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !isInitialized() {
			yellowBold("Please make sure you do 'forta init' first and check your configuration at %s/config.yml\n", cfg.FortaDir)
			return errors.New("not initialized")
		}
		return handler(cmd, args)
	}
}

func withDevOnly(handler func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !cfg.Development {
			return nil // no-op feature if not a dev run
		}
		return handler(cmd, args)
	}
}
