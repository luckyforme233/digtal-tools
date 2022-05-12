package cmd

import (
	"digtal/src/console"
	"digtal/src/digtal"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var list = &cobra.Command{
	Use:   "list",
	Short: "显示列表",
	Long:  `一个开机脚本`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		list, _ := digtal.DropletList()
		title := []string{
			"ID", "名字", "配置（cpu/ram/disk)", "区域/ip", "创建时间",
		}
		data := make([][]string, 0)
		for _, item := range list {
			data = append(data, []string{
				strconv.Itoa(item.ID),
				item.Name,
				fmt.Sprintf("%d/%d/%d", item.Size.Vcpus, item.Size.Memory, item.Size.Disk),
				fmt.Sprintf("%s/%s/%s", item.Region.Name, item.Networks.V4[0].IPAddress, item.Networks.V4[1].IPAddress),
				item.Created,
			})
		}
		console.Render(title, data)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		digtal.InitClient()
	},
}
