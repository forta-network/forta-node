package cmd

import (
	"errors"
	"fmt"
	"forta-network/forta-node/config"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyFortaDir        = "forta_dir"
	keyFortaConfigFile = "forta_config_file"
	keyFortaPassphrase = "forta_passphrase"
	keyFortaProduction = "forta_production"
)

var (
	cfg config.Config

	parsedArgs struct {
		PrivateKeyFilePath string
		AgentVersion       uint64
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
	}

	cmdFortaAgentAdd = &cobra.Command{
		Use:   "add",
		Short: "add new agent",
		RunE:  withInitialized(withValidConfig(handleFortaAgentAdd)),
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

	// Global (persistent) flags

	cmdForta.PersistentFlags().String("dir", "", "Forta dir (default is $HOME/.forta) (overrides $FORTA_DIR)")
	viper.BindPFlag(keyFortaDir, cmdForta.PersistentFlags().Lookup("dir"))

	cmdForta.PersistentFlags().String("config", "", "config file (default is $HOME/.forta/config.yml) (overrides $FORTA_CONFIG_FILE)")
	viper.BindPFlag(keyFortaConfigFile, cmdForta.PersistentFlags().Lookup("config"))

	cmdForta.PersistentFlags().Bool("production", false, "production mode (overrides $FORTA_PRODUCTION)")
	viper.BindPFlag(keyFortaProduction, cmdForta.PersistentFlags().Lookup("production"))

	cmdForta.PersistentFlags().String("passphrase", "", "passphrase to decrypt the private key (overrides $FORTA_PASSPHRASE)")
	viper.BindPFlag(keyFortaPassphrase, cmdForta.PersistentFlags().Lookup("passphrase"))

	// forta account import

	cmdFortaAccountImport.Flags().String("file", "", "path to a file that contains a private key hex")
	cmdFortaAccountImport.MarkFlagRequired("file")
}

func initConfig() {
	viper.SetConfigType("yaml")

	viper.BindEnv(keyFortaDir)
	viper.BindEnv(keyFortaConfigFile)
	viper.BindEnv(keyFortaPassphrase)
	viper.BindEnv(keyFortaProduction)
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
	cfg.Production = viper.GetBool(keyFortaProduction)
	cfg.Passphrase = viper.GetString(keyFortaPassphrase)
	cfg.LocalAgentsPath = path.Join(cfg.FortaDir, config.DefaultLocalAgentsFileName)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("failed to read the config file: %v", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Errorf("failed to unmarshal the config file: %v", err)
	}

	if err := config.InitLogLevel(cfg); err != nil {
		log.Errorf("failed to init log level: %v", err)
	}
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
		fmt.Println("The config file has invalid or missing fields:")
		for _, validationErr := range validationErrs {
			fmt.Printf("  - %s\n", validationErr.Namespace()[7:])
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
