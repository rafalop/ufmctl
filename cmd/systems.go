package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"os"
	"strconv"
)

var systemsCmd = &cobra.Command{
	Use:   "systems {subcommand} {args}",
	Short: "Perform systems operations",
	Long:  "Get systems in ufm",
	//Run: func(cmd *cobra.Command, args []string) {
	//	if len(args) < 1 {
	//		fmt.Println("pkeys requires at least one subcommand.")
	//	}
	//},
}

var systemsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List systems in UFM",
	//Long:  "List sytsems in UFM",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		filters := []string{}
		filters = append(filters, "brief="+strconv.FormatBool(SystemsBrief))
		filters = append(filters, "chassis="+strconv.FormatBool(SystemsChassis))
		filters = append(filters, "ports="+strconv.FormatBool(SystemsPorts))
		filters = append(filters, "in_rack="+strconv.FormatBool(SystemsInRack))
		if len(SystemsIP) > 0 {
			filters = append(filters, "ip="+SystemsIP)
		}
		if len(SystemsType) > 0 {
			filters = append(filters, "type="+SystemsType)
		}
		if len(SystemsModel) > 0 {
			filters = append(filters, "model="+SystemsModel)
		}
		if len(SystemsRole) > 0 {
			filters = append(filters, "role="+SystemsRole)
		}
		if len(SystemsPeerName) > 0 {
			filters = append(filters, "peer_name="+SystemsPeerName)
		}
		if len(SystemsComputes) > 0 {
			filters = append(filters, "computes="+SystemsComputes)
		}
		systems, err := u.GetSystems(filters)
		if err != nil {
			ExitError(err)
		}
		//systemsJson, err := json.Marshal(systems)
		//fmt.Println(string(systemsJson))
		if Format == "json" {
			fmt.Println(string(systems))
		} else {
			printSystemsTable(systems)
		}
	},
}

func printSystemsTable(systems string) {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	rowConfigAutomerge := table.RowConfig{AutoMerge: true}
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"STATE", "ROLE", "TYPE", "FW_VERSION", "IP", "GUID", "SYS_NAME", "PORTS"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 7, AutoMerge: true},
	})
	for _, system := range gjson.Get(systems, "#(*)#").Array() {
		t.AppendRow(table.Row{system.Get("state").String(), system.Get("role").String(), system.Get("type").String(), system.Get("fw_version").String(), system.Get("ip").String(), system.Get("guid"), system.Get("system_name").String()}, rowConfigAutomerge)
	}
	fmt.Println(t.Render())
}
