package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version string

var rootCmd = &cobra.Command{
	Use:   "netfetch",
	Short: "Netfetch is a CLI tool for scanning Kubernetes clusters for network policies",
	Long: `Netfetch is a CLI  tool for scanning clusters for network policies and identifying unprotected workloads.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Netfetch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(Version + "\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
