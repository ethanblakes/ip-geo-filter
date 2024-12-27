package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示ip-geo-filter的版本号",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0")
	},
}
