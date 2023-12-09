package cmd

import (
	"fmt"

	"github.com/deggja/netfetch/backend/pkg/k8s"
	"github.com/spf13/cobra"
)

var dryRun bool

var scanCmd = &cobra.Command{
	Use:   "scan [namespace]",
	Short: "Scan Kubernetes namespaces for network policies",
	Long: `Scan Kubernetes namespaces for network policies. 
	You can perform a dry run of the scan using the --dryrun or -d flag, 
	which will simulate the scan without making any changes.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var namespace string
		if len(args) > 0 {
			namespace = args[0]
		}
		_, err := k8s.ScanNetworkPolicies(namespace, dryRun, false, true, true, true)
		if err != nil {
			// Handle the error appropriately
			fmt.Println("Error during scan:", err)
			return
		}
	},
}

func init() {
	scanCmd.Flags().BoolVarP(&dryRun, "dryrun", "d", false, "Perform a dry run without applying any changes")
	rootCmd.AddCommand(scanCmd)
}
