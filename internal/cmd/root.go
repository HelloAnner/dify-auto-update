package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dify-auto-update",
	Short: "A tool to automatically sync local folders with Dify knowledge bases",
	Long: `dify-auto-update is a CLI tool that watches specified folders and automatically
syncs their contents with Dify knowledge bases. It maintains the same structure
between local folders and Dify knowledge bases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startWatcher()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// 配置文件已在 main.go 中加载
} 