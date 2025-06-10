package cmd

import (
	"fmt"
	//"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	//"github.com/tidwall/gjson"
	"os"
)

var actionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Perform an action on the UFM api server",
}

var actionsGetCablesInfoCmd = &cobra.Command{
	Use:   "get-cables-info [port id]",
	Args: cobra.ExactArgs(1),	
	Short: "Get detailed cable and module info for a switch port (json output)",
	Long: "port id must be the ufm port id, usually of the form 'guid_x'",
	Run: func(cmd *cobra.Command, args []string) {
		u := GetUfmClient()
		data, err := u.GetCablesInfo(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println(data)
		}
		os.Exit(0)
	},
}
