package cmd

import (
	"fmt"

	"github.com/deggja/netfetch/backend/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	dryRun bool
	native bool
	cilium bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [namespace]",
	Short: "Scan Kubernetes namespaces for network policies",
	Long: `Scan Kubernetes namespaces for network policies.
	You can specify --native or --cilium to scan only native or Cilium network policies respectively.
	Combining both or using none will scan for both types.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var namespace string
		if len(args) > 0 {
			namespace = args[0]
		}

		// Perform Kubernetes native network policy scan if no specific type is mentioned or if --native is used
		if !cilium || native {
			nativeScanResult, err := k8s.ScanNetworkPolicies(namespace, dryRun, false, true, true, true)
			if err != nil {
				fmt.Println("Error during kubernetes native network policies scan:", err)
			} else {
				fmt.Println("Kubernetes native network policies scan completed successfully.")
				// Do something with nativeScanResult like logging or aggregating data
				handleScanResult(nativeScanResult)
			}
		}

		// Perform Cilium network policy scan if no specific type is mentioned or if --cilium is used
		if !native || cilium {
			ciliumScanResult, err := k8s.ScanCiliumNetworkPolicies(namespace, dryRun, false, true, true, true)
			if err != nil {
				fmt.Println("Error during cilium network policies scan:", err)
			} else {
				fmt.Println("Cilium network policies scan completed successfully.")
				// Do something with ciliumScanResult like logging or aggregating data
				handleScanResult(ciliumScanResult)
			}
		}
	},
}

func handleScanResult(scanResult *k8s.ScanResult) {
	// Implement your logic to handle scan results
}

func init() {
	scanCmd.Flags().BoolVarP(&dryRun, "dryrun", "d", false, "Perform a dry run without applying any changes")
	scanCmd.Flags().BoolVar(&native, "native", false, "Scan only native network policies")
	scanCmd.Flags().BoolVar(&cilium, "cilium", false, "Scan only Cilium network policies")
	rootCmd.AddCommand(scanCmd)
}
