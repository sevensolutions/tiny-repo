package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tinyrepo",
	Short: "A very basic repository for versioned artifacts.",
	Long:  `A very basic repository for versioned artifacts.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	// helper.HandleError(viper.BindEnv("API_KEY"))
	// helper.HandleError(viper.BindEnv("API_SECRET"))
	// helper.HandleError(viper.BindEnv("USERNAME"))
	// helper.HandleError(viper.BindEnv("PASSWORD"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configuration file: ", viper.ConfigFileUsed())
	}
}
