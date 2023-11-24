package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Helper function to write to both buffer and standard output
func printToBoth(writer *bufio.Writer, s string) {
	// Print to standard output
	fmt.Print(s)
	// Write the same output to buffer
	fmt.Fprint(writer, s)
}

// ScanNetworkPolicies scans namespaces for network policies
func ScanNetworkPolicies(specificNamespace string) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %s\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %s\n", err)
		return
	}

	var namespacesToScan []string
	if specificNamespace != "" {
		namespacesToScan = append(namespacesToScan, specificNamespace)
	} else {
		allNamespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing namespaces: %s\n", err)
			return
		}
		for _, ns := range allNamespaces.Items {
			if !isSystemNamespace(ns.Name) {
				namespacesToScan = append(namespacesToScan, ns.Name)
			}
		}
	}

	missingPoliciesOrUncoveredPods := false
	userDeniedPolicyApplication := false
	deniedNamespaces := []string{}

	for _, nsName := range namespacesToScan {
		policies, err := clientset.NetworkingV1().NetworkPolicies(nsName).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorMsg := fmt.Sprintf("\nError listing network policies in namespace %s: %s\n", nsName, err)
			printToBoth(writer, errorMsg)
			continue
		}

		hasDenyAll := hasDefaultDenyAllPolicy(policies.Items)
		coveredPods := make(map[string]bool)
		hasPolicies := len(policies.Items) > 0

		for _, policy := range policies.Items {
			selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.PodSelector)
			if err != nil {
				fmt.Printf("Error parsing selector for policy %s: %s\n", policy.Name, err)
				continue
			}

			pods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{
				LabelSelector: selector.String(),
			})
			if err != nil {
				fmt.Printf("Error listing pods for policy %s: %s\n", policy.Name, err)
				continue
			}

			for _, pod := range pods.Items {
				coveredPods[pod.Name] = true
			}
		}

		var tableOutput strings.Builder
		if !hasPolicies || hasDenyAll {
			allPods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				errorMsg := fmt.Sprintf("Error listing all pods in namespace %s: %s\n", nsName, err)
				printToBoth(writer, errorMsg)
				continue
			}

			unprotectedPods := false
			var unprotectedPodDetails []string
			for _, pod := range allPods.Items {
				if !coveredPods[pod.Name] {
					missingPoliciesOrUncoveredPods = true
					unprotectedPods = true
					podDetail := fmt.Sprintf("%-30s %-20s %-15s", nsName, pod.Name, pod.Status.PodIP)
					unprotectedPodDetails = append(unprotectedPodDetails, podDetail)
				}
			}

			if unprotectedPods {
				if len(unprotectedPodDetails) > 0 {
					tableOutput.WriteString("\nNetfetch found the following unprotected pods:\n\n")
					tableOutput.WriteString(fmt.Sprintf("%-30s %-20s %-15s\n", "Namespace", "Pod Name", "Pod IP"))
					for _, detail := range unprotectedPodDetails {
						tableOutput.WriteString(detail + "\n")
					}
					fmt.Print(tableOutput.String())
					writer.WriteString(tableOutput.String())
				}

				confirm := false
				prompt := &survey.Confirm{
					Message: fmt.Sprintf("Do you want to add a default deny all network policy to the namespace %s?", nsName),
				}
				survey.AskOne(prompt, &confirm, nil)

				if confirm {
					err := createAndApplyDefaultDenyPolicy(clientset, nsName)
					if err != nil {
						errorPolicyMsg := fmt.Sprintf("\nFailed to apply default deny policy in namespace %s: %s\n", nsName, err)
						printToBoth(writer, errorPolicyMsg)
					} else {
						successPolicyMsg := fmt.Sprintf("\nApplied default deny policy in namespace %s\n", nsName)
						printToBoth(writer, successPolicyMsg)
					}
				} else {
					userDeniedPolicyApplication = true
					deniedNamespaces = append(deniedNamespaces, nsName)
				}
			}
		}
	}

	writer.Flush()
	if output.Len() > 0 {
		saveToFile := false
		prompt := &survey.Confirm{
			Message: "Do you want to save the output to netfetch.txt?",
		}
		survey.AskOne(prompt, &saveToFile, nil)

		if saveToFile {
			err := os.WriteFile("netfetch.txt", output.Bytes(), 0644)
			if err != nil {
				errorFileMsg := fmt.Sprintf("Error writing to file: %s\n", err)
				printToBoth(writer, errorFileMsg)
			} else {
				printToBoth(writer, "Output file created: netfetch.txt\n")
			}
		} else {
			printToBoth(writer, "Output file not created.\n")
		}
	}

	if missingPoliciesOrUncoveredPods {
		if userDeniedPolicyApplication {
			printToBoth(writer, "\nFor the following namespaces, you should assess the need of implementing network policies:\n")
			for _, ns := range deniedNamespaces {
				fmt.Println(" -", ns)
			}
			printToBoth(writer, "\nConsider either an implicit default deny all network policy or a policy that targets the pods not selected by a network policy. Check the Kubernetes documentation for more information on network policies: https://kubernetes.io/docs/concepts/services-networking/network-policies/\n")
		} else {
			printToBoth(writer, "\nNetfetch scan completed!\n")
		}
	} else {
		printToBoth(writer, "\nNo network policies missing. You are good to go!\n")
	}
}

// Function to create the implicit default deny if missing
func createAndApplyDefaultDenyPolicy(clientset *kubernetes.Clientset, namespace string) error {
	policyName := namespace + "-default-deny-all"
	policy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyName,
			Namespace: namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{}, // Selects all pods in the namespace
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
		},
	}

	_, err := clientset.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), policy, metav1.CreateOptions{})
	return err
}

// hasDefaultDenyAllPolicy checks if the list of policies includes a default deny all policy
func hasDefaultDenyAllPolicy(policies []networkingv1.NetworkPolicy) bool {
	for _, policy := range policies {
		if isDefaultDenyAllPolicy(policy) {
			return true
		}
	}
	return false
}

// isDefaultDenyAllPolicy checks if a single network policy is a default deny all policy
func isDefaultDenyAllPolicy(policy networkingv1.NetworkPolicy) bool {
	return len(policy.Spec.Ingress) == 0 && len(policy.Spec.Egress) == 0
}

// isSystemNamespace checks if the given namespace is a system namespace
func isSystemNamespace(namespace string) bool {
	switch namespace {
	case "kube-system", "tigera-operator", "kube-public", "kube-node-lease", "gatekeeper-system", "calico-system":
		return true
	default:
		return false
	}
}
