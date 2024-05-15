package cmd

import (
	"fmt"

	"github.com/deggja/netfetch/backend/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	dryRun       bool
	native       bool
	cilium       bool
	verbose      bool
	targetPolicy string
)

var scanCmd = &cobra.Command{
	Use:   "scan [namespace]",
	Short: "Scan Kubernetes namespaces for network policies",
	Long: `Scan Kubernetes namespaces for network policies.
    By default, it scans for native Kubernetes network policies.
    Use --cilium to scan for Cilium network policies. You may also target a speecific network policy using --target-policy.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var namespace string
		if len(args) > 0 {
			namespace = args[0]
		}

		// Initialize the Kubernetes clients
		clientset, err := k8s.GetClientset()
		if err != nil {
			fmt.Println("Error creating Kubernetes client:", err)
			return
		}
		dynamicClient, err := k8s.GetCiliumDynamicClient()
		if err != nil {
			fmt.Println("Error creating Kubernetes dynamic client:", err)
			return
		}

		// Handle target policy for native Kubernetes network policies
		if targetPolicy != "" {
			if !cilium || native {
				fmt.Println("Policy type: Kubernetes")
				fmt.Printf("Searching for Kubernetes native network policy '%s' across all non-system namespaces...\n", targetPolicy)
				policy, foundNamespace, err := k8s.FindNativeNetworkPolicyByName(dynamicClient, clientset, targetPolicy)
				if err != nil {
					fmt.Println("Error during Kubernetes native network policy search:", err)
				} else {
					fmt.Printf("Found Kubernetes native network policy '%s' in namespace '%s'.\n", policy.GetName(), foundNamespace)

					// List the pods targeted by this policy
					pods, err := k8s.ListPodsTargetedByNetworkPolicy(dynamicClient, policy, foundNamespace)
					if err != nil {
						fmt.Printf("Error listing pods targeted by policy %s: %v\n", policy.GetName(), err)
					} else if len(pods) == 0 {
						fmt.Printf("No pods targeted by policy '%s' in namespace '%s'.\n", policy.GetName(), foundNamespace)
					} else {
						fmt.Printf("Pods targeted by policy '%s' in namespace '%s':\n", policy.GetName(), foundNamespace)
						for _, pod := range pods {
							fmt.Printf("  - %s\n", pod)
						}
					}
				}
				return
			}
		}

		// Handle target policy for Cilium network policies and cluster-wide policies
        if targetPolicy != "" && cilium {
            fmt.Println("Policy type: Cilium")
            fmt.Printf("Searching for Cilium network policy '%s' across all non-system namespaces...\n", targetPolicy)
            policy, foundNamespace, err := k8s.FindCiliumNetworkPolicyByName(dynamicClient, targetPolicy)
            if err != nil {
                // If not found in namespaces, search for cluster-wide policy
                fmt.Println("Cilium network policy not found in namespaces, searching for cluster-wide policy...")
                policy, err = k8s.FindCiliumClusterWideNetworkPolicyByName(dynamicClient, targetPolicy)
                if err != nil {
                    fmt.Println("Error during Cilium cluster-wide network policy search:", err)
                } else {
                    fmt.Printf("Found Cilium cluster-wide network policy '%s'.\n", policy.GetName())

                    // List the pods targeted by this cluster-wide policy
                    pods, err := k8s.ListPodsTargetedByCiliumClusterWideNetworkPolicy(dynamicClient, policy)
                    if err != nil {
                        fmt.Printf("Error listing pods targeted by cluster-wide policy %s: %v\n", policy.GetName(), err)
                    } else if len(pods) == 0 {
                        fmt.Printf("No pods targeted by cluster-wide policy '%s'.\n", policy.GetName())
                    } else {
                        fmt.Printf("Pods targeted by cluster-wide policy '%s':\n", policy.GetName())
                        for _, pod := range pods {
                            fmt.Printf("  - %s\n", pod)
                        }
                    }
                }
            } else {
                fmt.Printf("Found Cilium network policy '%s' in namespace '%s'.\n", policy.GetName(), foundNamespace)

                // List the pods targeted by this policy
                pods, err := k8s.ListPodsTargetedByCiliumNetworkPolicy(dynamicClient, policy, foundNamespace)
                if err != nil {
                    fmt.Printf("Error listing pods targeted by policy %s: %v\n", policy.GetName(), err)
                } else if len(pods) == 0 {
                    fmt.Printf("No pods targeted by policy '%s' in namespace '%s'.\n", policy.GetName(), foundNamespace)
                } else {
                    fmt.Printf("Pods targeted by policy '%s' in namespace '%s':\n", policy.GetName(), foundNamespace)
                    for _, pod := range pods {
                        fmt.Printf("  - %s\n", pod)
                    }
                }
            }
            return
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
	scanCmd.Flags().StringVarP(&targetPolicy, "target", "t", "", "Scan a specific network policy by name")
	scanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.AddCommand(scanCmd)
}
