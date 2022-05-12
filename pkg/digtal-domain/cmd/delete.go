package cmd

import (
	"digtal/src/digtal"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var delete = &cobra.Command{
	Use:   "delete",
	Short: "显示列表",
	Long:  `一个开机脚本`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println("id 错误", err)
			return
		}
		digtal.DeleteDroplet(id)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		digtal.InitClient()
	},
}
