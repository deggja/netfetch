package cmd

import (
	"fmt"

	"github.com/deggja/netfetch/backend/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	dryRun  bool
	native  bool
	cilium  bool
	verbose bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [namespace]",
	Short: "Scan Kubernetes namespaces for network policies",
	Long: `Scan Kubernetes namespaces for network policies.
    By default, it scans for native Kubernetes network policies.
    Use --cilium to scan for Cilium network policies.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var namespace string
		if len(args) > 0 {
			namespace = args[0]
		}

		// Default to native scan if no specific type is mentioned or if --native is used
		if !cilium || native {
			fmt.Println("Running native network policies scan...")
			nativeScanResult, err := k8s.ScanNetworkPolicies(namespace, dryRun, false, true, true, true)
			if err != nil {
				fmt.Println("Error during Kubernetes native network policies scan:", err)
			} else {
				fmt.Println("Kubernetes native network policies scan completed successfully.")
				handleScanResult(nativeScanResult)
			}
		}

		// Perform Cilium network policy scan if --cilium is used
		if cilium {
			// Perform cluster-wide Cilium scan first if no namespace is specified
			if namespace == "" {
				fmt.Println("Running cluster-wide Cilium network policies scan...")
				dynamicClient, err := k8s.GetCiliumDynamicClient()
				if err != nil {
					fmt.Println("Error obtaining dynamic client:", err)
					return
				}

				clusterwideScanResult, err := k8s.ScanCiliumClusterwideNetworkPolicies(dynamicClient, false, dryRun, true)
				if err != nil {
					fmt.Println("Error during cluster-wide Cilium network policies scan:", err)
				} else {
					// Handle the cluster-wide scan result; skip further scanning if all pods are protected
					if clusterwideScanResult.AllPodsProtected {
						fmt.Println("All pods are protected by cluster wide cilium policies.\nYour Netfetch security score is: 42/42")
						return
					}
					handleScanResult(clusterwideScanResult)
				}
			}

			// Proceed with normal Cilium network policy scan
			fmt.Println("Running Cilium network policies scan for namespaces...")
			ciliumScanResult, err := k8s.ScanCiliumNetworkPolicies(namespace, dryRun, false, true, true, true)
			if err != nil {
				fmt.Println("Error during Cilium network policies scan:", err)
			} else {
				fmt.Println("Cilium network policies scan completed successfully.")
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
	scanCmd.Flags().BoolVar(&cilium, "cilium", false, "Scan only Cilium network policies (includes cluster-wide policies if no namespace is specified)")
	scanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.AddCommand(scanCmd)
}
