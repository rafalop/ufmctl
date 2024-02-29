package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"ufmctl/pkg/ufm"
)

var rootCmd = &cobra.Command{
	Use:   "ufmctl",
	Short: "cli for interacting with UFM api",
}

func Execute() {
	Init()
	rootCmd.Execute()
}

func GetUfmClient() *ufm.UfmClient {

	var err error
	UfmClient, err = ufm.GetClient(Username, Password, Endpoint, Insecure, CookieFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return UfmClient
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

var Username string
var Password string
var Endpoint string
var Insecure bool
var PrintStatus bool
var CookieFile string
var UfmClient *ufm.UfmClient
var Format string

// Pkeys flags
var (
	PkKeysOnly   bool
	PkIndex0     bool
	PkIpoIb      bool
	PkMembership string
)

var (
	PortsFilters     string
	PortsOutputBrief bool
)

// Systems filters
var (
	SystemsBrief    bool
	SystemsIP       string
	SystemsType     string
	SystemsModel    string
	SystemsRole     string
	SystemsPeerName string
	SystemsChassis  bool
	SystemsPorts    bool
	SystemsInRack   bool
	SystemsComputes string
)

func Init() {
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "username to connect to UFM API with")
	rootCmd.PersistentFlags().StringVarP(&Format, "format", "f", "table", "output format (table, csv, json)")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "password to connect to UFM API with")
	rootCmd.PersistentFlags().StringVarP(&Endpoint, "endpoint", "e", "", "UFM API endpoint")
	rootCmd.PersistentFlags().BoolVarP(&Insecure, "insecure", "i", false, "use https without cert validation")
	rootCmd.PersistentFlags().BoolVarP(&PrintStatus, "status", "s", true, "print status to stderr")
	rootCmd.PersistentFlags().StringVarP(&CookieFile, "cookiefile", "c", "ufm-cookies.txt", "file to store cookies")
	rootCmd.MarkPersistentFlagRequired("username")
	rootCmd.MarkPersistentFlagRequired("password")
	rootCmd.MarkPersistentFlagRequired("endpoint")
	rootCmd.AddCommand(pkeysCmd)
	pkeysCmd.AddCommand(pkeysListCmd)
	pkeysListCmd.Flags().BoolVarP(&PkKeysOnly, "keys-only", "", false, "list only keys without guid info")

	pkeysCmd.AddCommand(pkeysAddCmd)
	//pkeysCmd.PersistentFlags().BoolVarP(&PkGuids, "guids", "", false, "include guid data for pkeys")
	//pkeysCmd.PersistentFlags().BoolVarP(&PkPorts, "ports", "", false, "include guid data for pkeys")
	//pkeysCmd.AddCommand(pkeysGetCmd)
	//pkeysCmd.AddCommand(pkeysCreateCmd)
	pkeysAddCmd.Flags().BoolVarP(&PkIndex0, "index0", "", true, "set index0 by default")
	pkeysAddCmd.Flags().BoolVarP(&PkIpoIb, "ipoib", "", true, "set ip over ib")
	pkeysAddCmd.Flags().StringVarP(&PkMembership, "membership", "", "limited", "type of membership (full or limited)")

	pkeysCmd.AddCommand(pkeysRemoveCmd)
	//pkeysCreateCmd.Flags().IntVarP(&PkMtuLimit, "mtu-limit", "", 4, "mtu limit")
	//pkeysCreateCmd.Flags().FloatVarP(&PkRateLimit, "rate-limit", "", "2.5", "members must have a higher rate than this to be allowed to connect")
	//pkeysCmd.AddCommand(pkeysMemberCmd)

	rootCmd.AddCommand(portsCmd)
	portsCmd.AddCommand(portsListCmd)
	portsListCmd.Flags().StringVarP(&PortsFilters, "filters", "", "", "comma delimited list of filters to try and apply to list query, eg. active=true,system=mycomputer,sys_type=Switch")
	portsListCmd.Flags().BoolVarP(&PortsOutputBrief, "output-brief", "", true, "only print brief output with limited fields. If this is false, only json is output")
	//portsListCmd.Flags().StringVarP(&PortsExtraColumns, "extra-columns", "", "",  "comma delimited list of extra columns to print in table mode")
	portsCmd.AddCommand(portsGetCmd)

	rootCmd.AddCommand(systemsCmd)
	//rootCmd.PersistentFlags().StringVarP(&SystemsIp, "ip", "", "ip address for system to get info for")
	//rootCmd.PersistentFlags().StringVarP(&SystemsIp, "name", "", "ip address for system to get info for")

	systemsCmd.AddCommand(systemsListCmd)
	systemsListCmd.Flags().BoolVarP(&SystemsBrief, "brief", "b", true, "Provides a brief response with essential information only")
	systemsListCmd.Flags().StringVarP(&SystemsIP, "ip", "", "", "System IP address")
	systemsListCmd.Flags().StringVar(&SystemsType, "type", "", "Specifies the type of system (switch/host/gateway/router)")
	systemsListCmd.Flags().StringVar(&SystemsModel, "model", "", "Specific model of a switch")
	systemsListCmd.Flags().StringVar(&SystemsRole, "role", "", "Specifies the role of the system (core/tor/endpoint)")
	systemsListCmd.Flags().StringVar(&SystemsPeerName, "peer-name", "", "List of peer devices (comma-separated)")
	systemsListCmd.Flags().BoolVar(&SystemsChassis, "chassis", false, "Specifies whether to provide detailed module descriptions or only module names")
	systemsListCmd.Flags().BoolVar(&SystemsPorts, "ports", false, "Specifies whether to provide detailed port descriptions or only port names")
	systemsListCmd.Flags().BoolVar(&SystemsInRack, "in-rack", false, "Specifies whether to get all systems that belong to a rack or those that do not belong to any rack")
	systemsCmd.Flags().StringVar(&SystemsComputes, "computes", "", "Specifies whether to retrieve systems that are allocated or free for logical servers")

	//maybe create the ufm object here get token or whatever

}

func PrintColumn(val string, padding int) {
	fmt.Printf("%-*s", padding, val)
}

func ExitError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
