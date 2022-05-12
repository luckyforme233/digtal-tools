package cmd

import (
	"digtal/src/digtal"
	"fmt"
	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "显示列表",
	Long:  `一个开机脚本`,
	Run: func(cmd *cobra.Command, args []string) {
		droplet, response, err := digtal.CreateDroplet()
		if err != nil {
			fmt.Println("err", err)
			return
		}
		fmt.Println(response.String())
		fmt.Println(droplet.ID)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		digtal.InitClient()
	},
}
