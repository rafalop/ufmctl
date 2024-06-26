package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"encoding/json"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tidwall/gjson"
	"os"
	//"errors"
)

var vPortsCmd = &cobra.Command{
	Use:   "vports",
	Short: "Get vport information",
	Long:  "Get vport information",
}

var vPortsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list vports",
	Long:  "Get list of all vports",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		var vPortsJson string
		var err error
		vPortsJson, err = u.VPortsGetAll(VPortsPhysport)
		if err != nil {
			ExitError(err)
		}
		if Format == "json" {
			fmt.Println(vPortsJson)
		} else {
			printVPortsTable(vPortsJson, Format)
		}
		os.Exit(0)
	},
}

func printVPortsTable(portsJson string, format string) {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	t.AppendHeader(table.Row{"NAME", "VPORT_GUID", "VPORT_STATE", "SYS_NAME"})
	for _, p := range gjson.Parse(portsJson).Array() {
		t.AppendRow(table.Row{p.Get("port_name").String(), p.Get("virtual_port_guid").String(), p.Get("virtual_port_state").String(), p.Get("system_name").String()})
	}
	if format == "csv" {
		fmt.Println(t.RenderCSV())
	} else {
		fmt.Println(t.Render())
	}
}
