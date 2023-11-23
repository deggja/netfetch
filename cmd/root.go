package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netfetch",
	Short: "Netfetch is a CLI tool for scanning Kubernetes clusters for network policies",
	Long: `Netfetch is a CLI tool that scans Kubernetes clusters for network policies
	and evaluates them against best practices. It helps in ensuring that your
	cluster's network configurations adhere to security standards.

	Usage:
	netfetch [command]

	Available Commands:
	scan        Scan Kubernetes namespaces for network policies
	help        Help about any command

	Flags:
	-h, --help   help for netfetch`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
}
