package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"encoding/json"
	"github.com/tidwall/gjson"
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Get port information",
	Long:  "Get port information",

}

var portsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list ports",
	Long:  "Get a liswt of ports",
	Run: func(cmd *cobra.Command, args []string) {
		u:=GetUfmClient()
		ports, err := u.GetPortsBrief()
		//ports, err := u.GetPortsFull()
		if err != nil {
			ExitError(err)
		}	
		if Format == "json" {
			jsonPorts, err := json.Marshal(ports)	
			if err != nil {
				ExitError(err)
			}
			fmt.Println(string(jsonPorts))
		} else {
			PrintColumn("NAME", 25)	
			PrintColumn("LG_STATE", 15)	
			PrintColumn("PHYS_STATE", 15)	
			PrintColumn("PATH", 50)	
			fmt.Printf("\n")	
			for _, port := range ports {
				PrintColumn(port.Name, 25)
				PrintColumn(port.LogicalState, 15)
				PrintColumn(port.PhysicalState, 15)
				PrintColumn(port.Path, 50)
				fmt.Printf("\n")
			}
			
		}
	},

}

var portsGetCmd = &cobra.Command{
	Use:   "get {port name}",
	Short: "get detailed info for a port",
	Long:  "get detailed infor for a port",
	Run: func(cmd *cobra.Command, args []string) {
		u:=GetUfmClient()
		portData, err := u.GetPort(args[0])
		if err != nil {
			ExitError(err)
		}	
		jsonPort, err := json.Marshal(portData)	
		if err != nil {
			ExitError(err)
		}
		if Format == "json" {
			fmt.Println(string(jsonPort))
		} else {
			result := gjson.Parse(string(jsonPort))	
			result.ForEach(func (key, value gjson.Result) bool{
				fmt.Printf("%25s : %s\n", key.String(), value.String())	
				return true
			})
		}
	},

}
