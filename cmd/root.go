package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	e "github.com/tinkershack/meteomunch/errors"
	"github.com/tinkershack/meteomunch/logger"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "munch",
	Short: "Munch meteo data",
	Long:  `Munch meteo data`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error(e.FATAL, "err", err)
		os.Exit(1)
	}
}

// log is shared within cmd package
var log *slog.Logger

func init() {
	log = logger.NewTag("cmd")
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./munch.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find working directory.
		pwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in the working directory with name "munch" (without extension).
		viper.AddConfigPath(pwd)
		viper.AddConfigPath(".")
		viper.SetConfigType("yml")
		viper.SetConfigName("munch")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn(e.FAIL, "err", "find config file")
		} else {
			log.Warn(e.FAIL, "err", err, "description", "read config file")
		}
	} else {
		log.Info("found config file", "path", viper.ConfigFileUsed())
	}
}
