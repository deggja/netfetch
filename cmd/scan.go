package cmd

import (
	"fmt"
	"netfetch/k8s"

	"github.com/spf13/cobra"
)

var namespace string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a Kubernetes namespace",
	Long:  `Scan a specified Kubernetes namespace for pod-to-pod network traffic.`,
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" {
			fmt.Println("Please specify a namespace")
			return
		}
		k8s.ScanNamespace(namespace)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Specify the namespace to scan")
}
