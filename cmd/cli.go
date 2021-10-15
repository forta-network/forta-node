package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/forta-protocol/forta-node/config"
	"github.com/mitchellh/mapstructure"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyFortaDir         = "forta_dir"
	keyFortaConfigFile  = "forta_config_file"
	keyFortaPassphrase  = "forta_passphrase"
	keyFortaDevelopment = "forta_development"
)

var (
	cfg config.Config

	parsedArgs struct {
		PrivateKeyFilePath string
		Version            uint64
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
		RunE:  withContractAddresses(withInitialized(withValidConfig(handleFortaRun))),
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
		Use:   "import",
		Short: "import new scanner account (removes the old one)",
		RunE:  withInitialized(handleFortaAccountImport),
	}

	cmdFortaAgent = &cobra.Command{
		Use:   "agent",
		Short: "agent management",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Hidden: true,
	}

	cmdFortaAgentAdd = &cobra.Command{
		Use:   "add",
		Short: "try an agent by adding it to the local list",
		RunE:  withAgentRegContractAddress(withDevOnly(withInitialized(withValidConfig(handleFortaAgentAdd)))),
	}

	cmdFortaImages = &cobra.Command{
		Use:   "images",
		Short: "list the Forta node container images",
		RunE:  handleFortaImages,
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

	cmdForta.AddCommand(cmdFortaAgent)
	cmdFortaAgent.AddCommand(cmdFortaAgentAdd)

	cmdForta.AddCommand(cmdFortaImages)

	// Global (persistent) flags

	cmdForta.PersistentFlags().String("dir", "", "Forta dir (default is $HOME/.forta) (overrides $FORTA_DIR)")
	viper.BindPFlag(keyFortaDir, cmdForta.PersistentFlags().Lookup("dir"))

	cmdForta.PersistentFlags().String("config", "", "config file (default is $HOME/.forta/config.yml) (overrides $FORTA_CONFIG_FILE)")
	viper.BindPFlag(keyFortaConfigFile, cmdForta.PersistentFlags().Lookup("config"))

	cmdForta.PersistentFlags().Bool("development", false, "development mode (overrides $FORTA_DEVELOPMENT)")
	viper.BindPFlag(keyFortaDevelopment, cmdForta.PersistentFlags().Lookup("development"))

	cmdForta.PersistentFlags().String("passphrase", "", "passphrase to decrypt the private key (overrides $FORTA_PASSPHRASE)")
	viper.BindPFlag(keyFortaPassphrase, cmdForta.PersistentFlags().Lookup("passphrase"))

	// forta account import

	cmdFortaAccountImport.Flags().String("file", "", "path to a file that contains a private key hex")
	cmdFortaAccountImport.MarkFlagRequired("file")

	// forta agent add
	cmdFortaAgentAdd.Flags().Uint64Var(&parsedArgs.Version, "version", 0, "agent version")
}

func initConfig() {
	viper.SetConfigType("yaml")

	viper.BindEnv(keyFortaDir)
	viper.BindEnv(keyFortaConfigFile)
	viper.BindEnv(keyFortaPassphrase)
	viper.BindEnv(keyFortaDevelopment)
	viper.AutomaticEnv()

	if cfg.FortaDir = viper.GetString(keyFortaDir); cfg.FortaDir == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfg.FortaDir = path.Join(home, ".forta")
	}

	cfg.ConfigPath = viper.GetString(keyFortaConfigFile)
	if cfg.ConfigPath == "" {
		cfg.ConfigPath = path.Join(cfg.FortaDir, "config.yml")
	}
	viper.SetConfigFile(cfg.ConfigPath)

	cfg.KeyDirPath = path.Join(cfg.FortaDir, config.DefaultKeysDirName)
	cfg.Development = !viper.GetBool(keyFortaDevelopment)
	cfg.Passphrase = viper.GetString(keyFortaPassphrase)
	cfg.LocalAgentsPath = path.Join(cfg.FortaDir, config.DefaultLocalAgentsFileName)
	cfg.LocalAgents, _ = readLocalAgents()

	// Inject our decode hook(s) to do extra stuff on the values before we unmarshal.
	decoderOpts := viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			resolveEnvVarsHookFunc,
		),
	)

	viper.ReadInConfig()
	if err := viper.Unmarshal(&cfg, decoderOpts); err != nil {
		fmt.Printf("failed to unmarshal the config: %v\n", err)
	}

	config.InitLogLevel(cfg)
}

var configEnvVarRegexp = regexp.MustCompile(`\$[A-Z0-9_]+`)

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
			yellowBold("Please make sure you do 'forta init' first and check your configuration at %s\n", cfg.ConfigPath)
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

func withContractAddresses(handler func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := useEnsDefaults(); err != nil {
			return err
		}
		return handler(cmd, args)
	}
}

func withAgentRegContractAddress(handler func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := useEnsAgentReg(); err != nil {
			return err
		}
		return handler(cmd, args)
	}
}
