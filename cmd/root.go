package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "signature-service",
	Short: "Signature Service - signing API",
	Long:  "Signature Service - signing API this microservice was created as a part of recruitment process",
}

// Execute command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (yaml/json)")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	viper.SetEnvPrefix("SIG") // np. SIG_PORT, SIG_DB_TYPE
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config.json")
		viper.AddConfigPath(".")
	}

}
