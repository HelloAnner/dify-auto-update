package cmd

import (
	"fmt"
	"time"

	"github.com/HelloAnner/dify-auto-update/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start watching folder for changes",
	Long:  `Start watching the specified folder for changes and sync with Dify knowledge bases`,
	Run: func(cmd *cobra.Command, args []string) {
		folder := viper.GetString("folder")
		apiKey := viper.GetString("api-key")
		baseURL := viper.GetString("base-url")

		syncer := service.NewDifySyncer(baseURL, apiKey)
		watcher := service.NewFolderWatcher(folder, syncer)

		fmt.Printf("开始监听文件夹: %s\n", folder)
		fmt.Printf("使用Dify基础URL: %s\n", baseURL)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// Initial sync
		if err := watcher.SyncFolder(); err != nil {
			fmt.Printf("初始化同步失败: %v\n", err)
		}

		for range ticker.C {
			if err := watcher.SyncFolder(); err != nil {
				fmt.Printf("同步失败: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
} 