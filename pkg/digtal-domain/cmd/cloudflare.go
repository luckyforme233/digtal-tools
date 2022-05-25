package cmd

import (
	"digtal/pkg/digtal-domain/config"
	"digtal/src/cloudflare"
	"digtal/src/digtal"
	"fmt"
	"github.com/spf13/cobra"
)

var cc = &cobra.Command{
	Use:   "cc",
	Short: "DNS 管理脚本",
	Long:  `一个开机脚本`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		dns, _ := cloudflare.CreateDns(config.C.CLDomain, "127.0.0.1")
		fmt.Println("添加域名：", dns)
		cloudflare.ShowAllDns(config.C.CLDomain)

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		digtal.InitClient()
	},
}
