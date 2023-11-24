package cmd

import (
	"github.com/deggja/netfetch/pkg/k8s"
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
		k8s.ScanNetworkPolicies(namespace)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
