package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	watchFolder string
	apiKey      string
	baseURL     string
)

var rootCmd = &cobra.Command{
	Use:   "dify-auto-update",
	Short: "A tool to automatically sync local folders with Dify knowledge bases",
	Long: `dify-auto-update is a CLI tool that watches specified folders and automatically
syncs their contents with Dify knowledge bases. It maintains the same structure
between local folders and Dify knowledge bases.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dify-auto-update.yaml)")
	rootCmd.PersistentFlags().StringVar(&watchFolder, "folder", "", "folder to watch for changes")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Dify API key")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "http://192.168.101.236", "Dify base URL")

	rootCmd.MarkPersistentFlagRequired("folder")
	rootCmd.MarkPersistentFlagRequired("api-key")

	viper.BindPFlag("folder", rootCmd.PersistentFlags().Lookup("folder"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("base-url", rootCmd.PersistentFlags().Lookup("base-url"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".dify-auto-update")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
} 