package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"encoding/json"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tidwall/gjson"
	"os"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Work with events",
}

var eventsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list events",
	Long:  "list all events",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		var eventsJson string
		var err error
		eventsJson, err = u.EventsGetAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if Format == "json" {
			fmt.Println(eventsJson)
		} else if Format == "table" {
			if DescriptionOnly {
				for _, p := range gjson.Parse(eventsJson).Array() {
					desc := p.Get("description").String()
					fmt.Println(desc)
				}
				os.Exit(0)
			}
			printEventsTable(eventsJson)
		}
		os.Exit(0)
	},
}

var eventsGetCmd = &cobra.Command{
	Use:   "get {event id}",
	Short: "get an event",
	Long:  "get detailed info about an event",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		var eventsJson string
		var err error
		eventsJson, err = u.EventsGet(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if Format == "json" {
			fmt.Println(eventsJson)
		} else {
			result := gjson.Parse(eventsJson).Array()
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

func printEventsTable(eventsJson string) {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	t.AppendHeader(table.Row{"ID", "TIMESTAMP", "SEVERITY", "OBJ_NAME", "ENTITY", "NAME", "DESC"})
	var entity string
	for _, p := range gjson.Parse(eventsJson).Array() {
		entity = EntityFromPath(p.Get("object_path").String())
		desc := p.Get("description").String()
		if len(desc) > 100 {
			desc = desc[0:100]
		}
		//t.AppendRow(table.Row{p.Get("name").String(), p.Get("guid"), p.Get("logical_state").String(), p.Get("physical_state").String(), p.Get("enabled_speed").String(), p.Get("path").String(), p.Get("peer_node_description").String()})
		t.AppendRow(table.Row{p.Get("id").String(), p.Get("timestamp"), p.Get("severity"), p.Get("object_name"), entity, p.Get("name").String(), desc})
	}
	fmt.Println(t.Render())

}

