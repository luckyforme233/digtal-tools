package console

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

// title:  []string{"Name", "Sign", "Rating"}
// data: data := [][]string{
//  {"A", "北京冬奥会 666", "100"},
//  {"A", "北京冬奥会真棒", "150"},
//  {"B", "Happy New Year 2022!", "200"},
//  {"B", "开工大吉！", "300"},
// }

func Render(title []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(title)
	// 合并第一列内容相同的单元格
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)

	table.AppendBulk(data)

	table.Render()
}
