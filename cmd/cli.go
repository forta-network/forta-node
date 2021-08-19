package cmd

import (
	"errors"
	"fmt"
	"forta-network/forta-node/config"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyFortaConfigDir  = "forta_config_dir"
	keyFortaConfigFile = "forta_config_file"
	keyFortaPassphrase = "forta_passphrase"
	keyFortaProduction = "forta_production"
)

var (
	cfg      config.Config
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
		RunE:  withValidConfig(handleFortaRun),
	}
)

// Execute executes the root command.
func Execute() error {
	cmdForta.AddCommand(cmdFortaInit)
	cmdForta.AddCommand(cmdFortaRun)
	return cmdForta.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global (persistent) flags

	cmdForta.PersistentFlags().String("dir", "", "Forta config dir (default is $HOME/.forta)")
	viper.BindPFlag(keyFortaConfigDir, cmdForta.PersistentFlags().Lookup("dir"))

	cmdForta.PersistentFlags().String("config", "", "config file (default is $HOME/.forta/config.yml)")
	viper.BindPFlag(keyFortaConfigFile, cmdForta.PersistentFlags().Lookup("config"))

	cmdForta.PersistentFlags().Bool("production", false, "production mode")
	viper.BindPFlag(keyFortaProduction, cmdForta.PersistentFlags().Lookup("production"))

	cmdForta.Flags().String("passphrase", "", "passphrase to decrypt the private key")
	viper.BindPFlag(keyFortaPassphrase, cmdForta.PersistentFlags().Lookup("passphrase"))
}

func initConfig() {
	viper.SetConfigType("yaml")

	viper.BindEnv(keyFortaConfigDir)
	viper.BindEnv(keyFortaConfigFile)
	viper.BindEnv(keyFortaPassphrase)
	viper.BindEnv(keyFortaProduction)
	viper.AutomaticEnv()

	if cfg.FortaDir = viper.GetString(keyFortaConfigDir); cfg.FortaDir == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfg.FortaDir = path.Join(home, ".forta")
	}

	cfg.ConfigPath = viper.GetString(keyFortaConfigFile)
	if cfg.ConfigPath == "" {
		cfg.ConfigPath = path.Join(cfg.FortaDir, "config.yml")
	}
	viper.SetConfigFile(cfg.ConfigPath)

	cfg.KeyDirPath = path.Join(cfg.FortaDir, ".keys")
	cfg.Production = viper.GetBool(keyFortaProduction)
	cfg.Passphrase = viper.GetString(keyFortaPassphrase)

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("failed to read the config file: %v", err))
	}

	if err := config.InitLogLevel(cfg); err != nil {
		log.Fatalf("failed to init log level: %v", err)
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
			fmt.Println("  -", validationErr.Namespace()[7:])
		}
		return errors.New("failed to validate the config file")
	}

	return nil
}

func withValidConfig(handler func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !isInitialized() {
			yellowBold("Please make sure you do 'forta init' first and check your configuration at %s\n", cfg.ConfigPath)
			return fmt.Errorf("not initialized")
		}
		if err := validateConfig(); err != nil {
			return err
		}
		return handler(cmd, args)
	}
}
