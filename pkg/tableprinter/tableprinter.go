package tableprinter

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// PrintTable prints a table to the console.
func PrintTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetCenterSeparator("|")
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
