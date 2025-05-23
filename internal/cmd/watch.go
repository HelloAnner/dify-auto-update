package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/HelloAnner/dify-auto-update/internal/service"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start watching folder for changes",
	Long:  `Start watching the specified folder for changes and sync with Dify knowledge bases`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := startWatcher(); err != nil {
			fmt.Printf("Error starting watcher: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func startWatcher() error {
	// 从配置文件读取配置
	watchFolder := viper.GetString("watch.folder")
	if watchFolder == "" {
		return fmt.Errorf("watch folder not configured in config.yaml")
	}

	interval := viper.GetDuration("watch.interval")
	if interval == 0 {
		interval = 5 * time.Minute // 默认5分钟
	} else {
		interval = interval * time.Second // 配置文件中的值以秒为单位
	}

	// 初始化 Dify 服务
	difyService := service.NewDifySyncer(
		viper.GetString("dify.base_url"),
		viper.GetString("dify.api_key"),
	)

	fmt.Println(difyService)

	// 创建文件夹监视器
	watcher := service.NewFolderWatcher(watchFolder, difyService)

	// 创建文件系统监控器
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer fsWatcher.Close()

	// 启动定时同步
	ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// 初始同步
		if err := watcher.SyncFolder(); err != nil {
			log.Printf("Initial sync error: %v", err)
		} else {
			log.Printf("Initial sync completed successfully")
		}

		for {
			select {
			case <-ticker.C:
				if err := watcher.SyncFolder(); err != nil {
				log.Printf("Error syncing folder: %v", err)
			}
		}
	}
}
