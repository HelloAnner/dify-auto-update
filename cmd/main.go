package main

import (
	"log"
	"os"

	"github.com/HelloAnner/dify-auto-update/internal/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	requiredConfigs := []string{
		"dify.api_key",
		"dify.base_url",
		"watch.folder",
	}

	for _, config := range requiredConfigs {
		if !viper.IsSet(config) {
			log.Printf("Required configuration '%s' is missing in config.yaml\n", config)
			os.Exit(1)
		}
	}

	log.Printf("Configuration loaded successfully from: %s", viper.ConfigFileUsed())

	if err := cmd.Execute(); err != nil {
		log.Printf("Error executing command: %s\n", err)
		os.Exit(1)
	}
} 