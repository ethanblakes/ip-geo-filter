package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "igf",
		Short: "一个用于IP地理位置查询的工具",            // 添加简短描述
		Long:  `本作品仅供学习参考使用，任何用其非法用途与作者无关。`, // 添加详细描述
	}
)

func init() {
	rootCmd.Version = "0.0.1"
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
