package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"encoding/json"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tidwall/gjson"
	"os"
)

var alarmsCmd = &cobra.Command{
	Use:   "alarms",
	Short: "Work with alarms",
}

var alarmsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list alarms",
	Long:  "list all alarms or all alarms for a device",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		var alarmsJson string
		var err error
		alarmsJson, err = u.AlarmsGetAll(AlarmsDeviceId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if Format == "json" {
			fmt.Println(alarmsJson)
		} else if Format == "table" {
			if DescriptionOnly {
				for _, p := range gjson.Parse(alarmsJson).Array() {
					desc := p.Get("reason").String()
					fmt.Println(desc)
				}
				os.Exit(0)
			}
			printAlarmsTable(alarmsJson)
		}
		os.Exit(0)
	},
}

var alarmsGetCmd = &cobra.Command{
	Use:   "get {alarm id}",
	Short: "get an alarm",
	Long:  "get detailed info about an alarm",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		var alarmsJson string
		var err error
		alarmsJson, err = u.AlarmsGet(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if Format == "json" {
			fmt.Println(alarmsJson)
		} else {
			result := gjson.Parse(alarmsJson).Array()
			if len(result) > 0 {
				result[0].ForEach(func(key, value gjson.Result) bool {
					fmt.Printf("%25s : %s\n", key.String(), value.String())
					return true
				})
			}
		}
		os.Exit(0)
	},
}

func printAlarmsTable(alarmsJson string) {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	t.AppendHeader(table.Row{"ID", "TIMESTAMP", "SEVERITY", "OBJ_NAME", "ENTITY", "REASON"})
	var entity string
	for _, p := range gjson.Parse(alarmsJson).Array() {
		entity = EntityFromPath(p.Get("object_path").String())
		t.AppendRow(table.Row{p.Get("id").String(), p.Get("timestamp"), p.Get("severity"), p.Get("object_name").String(), entity, p.Get("reason")})
	}
	fmt.Println(t.Render())

}

