package cmd

import (
	"fmt"

	"github.com/deggja/netfetch/backend/pkg/k8s"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [namespace]",
	Short: "Scan Kubernetes namespaces for network policies",
	Long:  `Scan all non-system Kubernetes namespaces for network policies and compare them with predefined standards.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var namespace string
		if len(args) > 0 {
			namespace = args[0]
		}
		_, err := k8s.ScanNetworkPolicies(namespace, false, true)
		if err != nil {
			// Handle the error appropriately
			fmt.Println("Error during scan:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
