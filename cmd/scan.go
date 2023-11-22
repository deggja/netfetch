package cmd

import (
	"github.com/deggja/netfetch/pkg/k8s"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan Kubernetes namespaces for network policies",
	Long:  `Scan all non-system Kubernetes namespaces for network policies and compare them with predefined standards.`,
	Run: func(cmd *cobra.Command, args []string) {
		k8s.ScanNetworkPolicies()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
