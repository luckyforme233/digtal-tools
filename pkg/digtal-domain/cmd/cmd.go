package cmd

import (
	"digtal/pkg/digtal-domain/config"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "digitalocean",
	Short: "digitalocean 开机脚本",
	Long:  `一个开机脚本`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	config.InitConfig()
	rootCmd.AddCommand(list)
	rootCmd.AddCommand(create)
	rootCmd.AddCommand(delete)
	rootCmd.AddCommand(one_key)
	rootCmd.AddCommand(cc)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
