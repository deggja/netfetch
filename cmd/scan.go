package cmd

import (
	"fmt"

	"github.com/deggja/netfetch/pkg/k8s"
	"github.com/deggja/netfetch/pkg/utils"
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

		podNames, err := k8s.ListPods(namespace)
		if err != nil {
			fmt.Printf("Error listing pods: %s\n", err)
			return
		}

		err = k8s.SniffPods(podNames, namespace)
		if err != nil {
			fmt.Printf("Error sniffing pods: %s\n", err)
			return
		}

		if !utils.CheckDependency("kubectl") {
			fmt.Println("kubectl is not installed. ", utils.InstallInstructions("kubectl"))
			return
		}
		if !utils.CheckDependency("ksniff") {
			fmt.Println("ksniff is not installed. ", utils.InstallInstructions("ksniff"))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Specify the namespace to scan")
}
