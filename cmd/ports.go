package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"encoding/json"
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tidwall/gjson"
	"os"
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Get port information",
	Long:  "Get port information",
}

var portsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list ports",
	Long:  "Get a list of ports",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		var portsJson string
		var err error
		portsJson, err = u.PortsGetAll(PortsFilters)
		if err != nil {
			ExitError(err)
		}
		if Format == "json" {
			fmt.Println(portsJson)
		} else {
			printPortsTable(portsJson, Format)
		}
		os.Exit(0)
	},
}

func printPortsTable(portsJson string, format string) {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	t.AppendHeader(table.Row{"NAME", "LG_STATE", "PHYS_STATE", "SYS_NAME", "NODE_DESC", "PATH"})
	for _, p := range gjson.Parse(portsJson).Array() {
		t.AppendRow(table.Row{p.Get("name").String(), p.Get("logical_state").String(), p.Get("physical_state").String(), p.Get("system_name").String(), p.Get("node_description").String(), p.Get("path").String()})
	}
	if format == "csv" {
		fmt.Println(t.RenderCSV())
	} else {
		fmt.Println(t.Render())
	}
}

var portsGetCmd = &cobra.Command{
	Use:   "get {port name}",
	Short: "get detailed info for a port",
	Long:  "get detailed infor for a port",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		portJson, err := u.PortsGet(args[0])
		if err != nil {
			ExitError(err)
		}
		if Format == "json" {
			fmt.Println(portJson)
		} else {
			result := gjson.Parse(portJson).Array()
			if len(result) > 0 {
				result[0].ForEach(func(key, value gjson.Result) bool {
					fmt.Printf("%25s : %s\n", key.String(), value.String())
					return true
				})
			}
		}
	},
}

var portsActionCmd = &cobra.Command{
	Use:   "action {enable|disable|reset} {port name}",
	Short: "action a port (disable, enable, reset)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "enable" && args[0] != "disable" && args[0] != "reset" {
			ExitError(errors.New(fmt.Sprintf("Incorrect action %s - you can only enable|disable|reset a port.", args[0])))
		}
		u := GetUfmClient()
		result, err := u.PortsAction(args[1], args[0])
		if err != nil {
			ExitError(err)
		}
		if Format == "json" {
			fmt.Println("Job created successfully: ", result)
		} else {
			fmt.Println("Job created successfully: ", result)
		}
	},
}
