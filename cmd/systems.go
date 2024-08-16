package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	//"os"
	"strconv"
	"strings"
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
var systemsGetCmd = &cobra.Command{
	Use:   "get {system name}",
	Short: "Get info for a system in UFM",
	//Long:  "List sytsems in UFM",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		systems, err := u.GetSystems([]string{})
		if err != nil {
			ExitError(err)
		}
		var systemID string
		for _, system := range gjson.Get(systems, "#(*)#").Array() {
			systemName := system.Get("system_name").String()
			if systemName == args[0] {
				systemID = system.Get("guid").String()
				PortsFilters = PortsFilters + `system=` + systemID
				break
			}
		}
		system, err := u.GetSystem(systemID)
		if err != nil {
			ExitError(err)
		}
		//systemsJson, err := json.Marshal(systems)
		//fmt.Println(string(systemsJson))
		if SystemsGuids {
			// get port list also to extract hca data
			//port, err := u.GetPorts
		} else {
			if Format == "json" {
				fmt.Println(string(system))
			} else {
				printSystemsTable(system)
			}
		}
	},
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
	//t.Style().Options.SeparateRows = true
	rowConfigAutomerge := table.RowConfig{AutoMerge: true}
	//t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"STATE", "ROLE", "TYPE", "FW_VERSION", "IP", "GUID", "SYS_NAME", "PORTS"})
	//t.SetColumnConfigs([]table.ColumnConfig{
	//	{Number: 1, AutoMerge: true},
	//	{Number: 2, AutoMerge: true},
	//	{Number: 3, AutoMerge: true},
	//	{Number: 4, AutoMerge: true},
	//	{Number: 5, AutoMerge: true},
	//	{Number: 6, AutoMerge: true},
	//	{Number: 7, AutoMerge: true},
	//	{Number: 8, AutoMerge: true},
	//})
	for _, system := range gjson.Get(systems, "#(*)#").Array() {
		ports := system.Get("ports")
		ports.ForEach(func(key, value gjson.Result) bool {
			printedPort := strings.Split(value.String(), "_")[0]
			//t.AppendRow(table.Row{system.Get("state").String(), system.Get("role").String(), system.Get("type").String(), system.Get("fw_version").String(), system.Get("ip").String(), system.Get("guid"), system.Get("system_name").String(), printedPort}, rowConfigAutomerge)
			t.AppendRow(table.Row{system.Get("state").String(), system.Get("role").String(), system.Get("type").String(), system.Get("fw_version").String(), system.Get("ip").String(), system.Get("guid"), system.Get("system_name").String(), printedPort}, rowConfigAutomerge)
			return true
		})
		//for _, port := range gjson.Get(system, ".ports#" {
		//	t.AppendRow(table.Row{system.Get("state").String(), system.Get("role").String(), system.Get("type").String(), system.Get("fw_version").String(), system.Get("ip").String(), system.Get("guid"), system.Get("system_name").String(), port.Get("0").String()}, rowConfigAutomerge)
		//}
		//t.AppendRow(table.Row{system.Get("state").String(), system.Get("role").String(), system.Get("type").String(), system.Get("fw_version").String(), system.Get("ip").String(), system.Get("guid"), system.Get("system_name").String()})
	}
	fmt.Println(t.Render())
}
