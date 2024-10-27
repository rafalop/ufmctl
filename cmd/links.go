package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"encoding/json"
	//"github.com/jedib0t/go-pretty/v6/table"
	//"github.com/tidwall/gjson"
	"os"
)

var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "Work with links",
}

var linksListCmd = &cobra.Command{
	Use:   "list",
	Short: "list links",
	Long:  "Get a list of links",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		var linksJson string
		var err error
		linksJson, err = u.LinksGetAll()
		if err != nil {
			ExitError(err)
		}
		//fmt.Println(portsJson)
		if Format == "json" {
			fmt.Println(linksJson)
		} else {
			//printPortsTable(portsJson, Format)
		}
		os.Exit(0)
	},
}

//func printPortsTable(portsJson string, format string) {
//	t := table.NewWriter()
//	t.Style().Options = table.OptionsNoBordersAndSeparators
//	t.AppendHeader(table.Row{"NAME", "GUID", "LG_STATE", "PHYS_STATE", "SPEEDS", "DESC", "PEER", "PEER_NODE"})
//	for _, p := range gjson.Parse(portsJson).Array() {
//		//t.AppendRow(table.Row{p.Get("name").String(), p.Get("guid"), p.Get("logical_state").String(), p.Get("physical_state").String(), p.Get("enabled_speed").String(), p.Get("path").String(), p.Get("peer_node_description").String()})
//		t.AppendRow(table.Row{p.Get("name").String(), p.Get("guid"), p.Get("logical_state").String(), p.Get("physical_state").String(), p.Get("enabled_speed").String(), p.Get("node_description"), p.Get("peer").String(), p.Get("peer_node_description").String()})
//	}
//	if format == "csv" {
//		fmt.Println(t.RenderCSV())
//	} else {
//		fmt.Println(t.Render())
//	}
//}
//
//var portsGetCmd = &cobra.Command{
//	Use:   "get {port name}",
//	Short: "get detailed info for a port",
//	Long:  "get detailed infor for a port",
//	Run: func(cmd *cobra.Command, args []string) {
//		u := GetUfmClient()
//		portJson, err := u.PortsGet(args[0])
//		if err != nil {
//			ExitError(err)
//		}
//		if Format == "json" {
//			fmt.Println(portJson)
//		} else {
//			result := gjson.Parse(portJson).Array()
//			if len(result) > 0 {
//				result[0].ForEach(func(key, value gjson.Result) bool {
//					fmt.Printf("%25s : %s\n", key.String(), value.String())
//					return true
//				})
//			}
//		}
//	},
//}
