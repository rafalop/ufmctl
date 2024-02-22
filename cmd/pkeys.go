package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/jedib0t/go-pretty/v6/table"
)

var pkeysCmd = &cobra.Command{
	Use:   "pkeys {subcommand} {args}",
	Short: "Perform pkey operations",
	Long:  "list, get or modify pkey settings",
	//Run: func(cmd *cobra.Command, args []string) {
	//	if len(args) < 1 {
	//		fmt.Println("pkeys requires at least one subcommand.")
	//	}	
	//},
}


var pkeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pkeys",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		pkeys, err := u.PkeyList()
		if err != nil {
			ExitError(err)
		}
		if PkKeysOnly {
			for pkey, _ := range gjson.Parse(pkeys).Map() {
				fmt.Println(pkey)
			}
			os.Exit(0)
		}
		if Format == "json" {
			fmt.Println(pkeys)
		}else if Format == "csv"{
			printPkeysTable(pkeys, "csv")
		} else { // assume table
			printPkeysTable(pkeys, "table")
		}
		os.Exit(0)
	},
}


func printPkeysTable(pkeys string, format string) {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	t.AppendHeader(table.Row{"PKEY", "GUID", "MSHIP", "INDEX0", "PORT_T", "IP", "PORT_NO", "DNAME", "IPOIB", "HOSTNAME", "NODE_DESC"})
	for pkey, rec := range gjson.Parse(pkeys).Map() {
		guidsArray := rec.Get("guids").Array()
		for _, g := range guidsArray {
			t.AppendRow(table.Row{pkey, g.Get("guid").String(), g.Get("membership").String(), g.Get("index0").String(), g.Get("port_type").String(), g.Get("ip").String(), g.Get("port_number").String(), g.Get("dname").String(), rec.Get("ip_over_ib").String(), g.Get("hostname").String(), g.Get("node_description").String() })
		}
	}
	if format == "csv" {
		fmt.Println(t.RenderCSV())
	} else {
		fmt.Println(t.Render())
	}
}

var pkeysAddCmd = &cobra.Command{
	Use:   "add {pkey} {guid1} {guid2} {guid3}...",
	Short: "Add members (GUIDs) to a pkey",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		fmt.Printf("addings guids to pkey=%s, guids=%v  ipoib=%t  index0=%t  membership=%s\n", args[0], args[1:], PkIpoIb, PkIndex0, PkMembership)
		err := u.PkeyAddGuids(args[0], PkIndex0, PkIpoIb, PkMembership, args[1:])
		if err != nil {
			ExitError(err)
		}
		os.Exit(0)
	},
}

var pkeysRemoveCmd = &cobra.Command{
	Use:   "remove {pkey} {guid1} {guid2} {guid3}...",
	Short: "Remove members (GUIDs) from pkey",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		fmt.Printf("removing guids from pkey=%s, guids=%v\n" ,args[0], args[1:])
		err := u.PkeyRemoveGuids(args[0], args[1:])
		if err != nil {
			ExitError(err)
		}
		os.Exit(0)
	},
}

//var pkeysCreateCmd = &cobra.Command{
//	Use:   "create {pkey}",
//	Short: "create a new pkey",
//	Run: func(cmd *cobra.Command, args []string) {
//		// command to connect to UFM goes here.
//		// process json
//		u:=GetUfmClient()
//		data := &ufm.CreatePkeyData{
//			Pkey: args[0],
//			Index0: PkIndex0,
//			IpOverIb: PkIpoib,
//			MtuLimit: PkMtulimit,
//			RateLimit: PkRateLimit,
//		}
//		err := u.CreatePkey(data)
//		//var queries = make([]string, 0)
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		fmt.Printf("Pkey %s created successfully.\n", args[0])
//	},
//}

//var pkeysMemberCmd = &cobra.Command{
//	Use:   "{pkey}",
//	Short: "add/remove members from a pkey",
//	Long: "update pkey memberships",
//	Run: func(cmd *cobra.Command, args []string) {
//		// command to connect to UFM goes here.
//		// process json
//		u:=GetUfmClient()
//		data := &ufm.CreatePkeyData{
//			Pkey: args[0],
//			Index0: PkIndex0,
//			IpOverIb: PkIpoib,
//			MtuLimit: PkMtulimit,
//			RateLimit: PkRateLimit,
//		}
//		err := u.CreatePkey(data)
//		//var queries = make([]string, 0)
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		fmt.Printf("Pkey %s created successfully.\n", args[0])
//	},
//}


