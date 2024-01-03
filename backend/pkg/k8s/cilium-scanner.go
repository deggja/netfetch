package k8s

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getCiliumDynamicClient returns a dynamic interface to query for Cilium policies
func getCiliumDynamicClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create k8s client config: %s", err)
		}
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return dynamicClient, nil
}

var hasStartedCiliumScan bool = false

// ScanCiliumNetworkPolicies scans namespaces for Cilium network policies
func ScanCiliumNetworkPolicies(specificNamespace string, dryRun bool, returnResult bool, isCLI bool, printScore bool, printMessages bool) (*ScanResult, error) {
	var output bytes.Buffer
	var namespacesToScan []string

	unprotectedPodsCount := 0
	scanResult := new(ScanResult)

	writer := bufio.NewWriter(&output)

	dynamicClient, err := getCiliumDynamicClient()
	if err != nil {
		fmt.Printf("Error creating dynamic Kubernetes client: %s\n", err)
		return nil, err
	}

	if dynamicClient == nil {
		fmt.Println("Failed to create dynamic client: client is nil")
		return nil, fmt.Errorf("failed to create dynamic client: client is nil")
	}

	clientset, err := GetClientset()
	if err != nil {
		fmt.Printf("Error creating Kubernetes clientset: %s\n", err)
		return nil, err
	}

	if clientset == nil {
		fmt.Println("Failed to create clientset: clientset is nil")
		return nil, fmt.Errorf("failed to create clientset: clientset is nil")
	}

	ciliumNPResource := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v2",
		Resource: "ciliumnetworkpolicies",
	}

	// Check if a specific namespace is provided
	if specificNamespace != "" {
		// Verify ns exists
		_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), specificNamespace, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				// Namespace does not exist
				return nil, fmt.Errorf("namespace %s does not exist", specificNamespace)
			}
			return nil, fmt.Errorf("error checking namespace %s: %v", specificNamespace, err)
		}
		namespacesToScan = append(namespacesToScan, specificNamespace)
	} else {
		namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("error listing namespaces: %v", err)
		}
		for _, ns := range namespaceList.Items {
			if !IsSystemNamespace(ns.Name) {
				namespacesToScan = append(namespacesToScan, ns.Name)
			}
		}
	}

	missingPoliciesOrUncoveredPods := false
	userDeniedPolicyApplication := false
	policyChangesMade := false
	deniedNamespaces := []string{}

	if isCLI && !hasStartedCiliumScan {
		fmt.Println("Policy type: Cilium")
		hasStartedCiliumScan = true
	}

	for _, nsName := range namespacesToScan {
		policies, err := dynamicClient.Resource(ciliumNPResource).Namespace(nsName).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorMsg := fmt.Sprintf("\nError listing Cilium network policies in namespace %s: %s\n", nsName, err)
			printToBoth(writer, errorMsg)
			return nil, errors.New(errorMsg)
		}

		var unstructuredPolicies []*unstructured.Unstructured
		for i := range policies.Items {
			unstructuredPolicies = append(unstructuredPolicies, &policies.Items[i])
		}

		hasDenyAll := hasDefaultDenyAllCiliumPolicy(unstructuredPolicies)
		coveredPods := make(map[string]bool)

		for _, policyUnstructured := range policies.Items {
			if isDefaultDenyAllCiliumPolicy(policyUnstructured) {
				hasDenyAll = true
			}
			policyMap := policyUnstructured.UnstructuredContent()

			spec, found := policyMap["spec"].(map[string]interface{})
			if !found {
				fmt.Fprintf(writer, "Error finding spec for policy %s in namespace %s\n", policyUnstructured.GetName(), nsName)
				continue
			}

			endpointSelector, found := spec["endpointSelector"].(map[string]interface{})
			if !found {
				fmt.Fprintf(writer, "Error finding endpointSelector for policy %s in namespace %s\n", policyUnstructured.GetName(), nsName)
				continue
			}

			labelSelector, err := convertEndpointToSelector(endpointSelector)
			if err != nil {
				fmt.Fprintf(writer, "Error converting endpoint selector to label selector for policy %s: %s\n", policyUnstructured.GetName(), err)
				continue
			}

			pods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{
				LabelSelector: labelSelector,
			})
			if err != nil {
				fmt.Fprintf(writer, "Error listing pods for endpointSelector %s: %s\n", labelSelector, err)
				continue
			}

			for _, pod := range pods.Items {
				coveredPods[pod.Name] = true
			}
		}

		if !hasDenyAll {
			var unprotectedPodDetails []string
			allPods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				errorMsg := fmt.Sprintf("Error listing all pods in namespace %s: %s\n", nsName, err)
				printToBoth(writer, errorMsg)
				continue
			}

			for _, pod := range allPods.Items {
				if !coveredPods[pod.Name] {
					podDetail := fmt.Sprintf("%s %s %s", nsName, pod.Name, pod.Status.PodIP)
					unprotectedPodDetails = append(unprotectedPodDetails, podDetail)
					unprotectedPodsCount++
				}
			}

			if len(unprotectedPodDetails) > 0 {
				missingPoliciesOrUncoveredPods = true
				scanResult.UnprotectedPods = append(scanResult.UnprotectedPods, unprotectedPodDetails...)
				// If CLI mode, interact with the user
				if isCLI {
					printToBoth(writer, "\nUnprotected Pods found in namespace "+nsName+":\n")
					for _, detail := range unprotectedPodDetails {
						printToBoth(writer, detail+"\n")
					}

					if !dryRun {
						confirm := false
						prompt := &survey.Confirm{
							Message: fmt.Sprintf("Do you want to add a default deny all Cilium network policy to the namespace %s?", nsName),
						}
						survey.AskOne(prompt, &confirm, nil)

						if confirm {
							err := createAndApplyDefaultDenyCiliumPolicy(nsName, dynamicClient)
							if err != nil {
								fmt.Printf("\nFailed to apply default deny Cilium policy in namespace %s: %s\n", nsName, err)
							} else {
								fmt.Printf("\nApplied default deny Cilium policy in namespace %s\n", nsName)
								policyChangesMade = true
							}
						} else {
							userDeniedPolicyApplication = true
							deniedNamespaces = append(deniedNamespaces, nsName)
						}
					}
				} else {
					// Non-CLI behavior
					scanResult.DeniedNamespaces = append(scanResult.DeniedNamespaces, nsName)
				}
			}
		}
	}

	writer.Flush()
	if output.Len() > 0 {
		saveToFile := false
		prompt := &survey.Confirm{
			Message: "Do you want to save the output to netfetch-cilium.txt?",
		}
		survey.AskOne(prompt, &saveToFile, nil)

		if saveToFile {
			err := os.WriteFile("netfetch-cilium.txt", output.Bytes(), 0644)
			if err != nil {
				errorFileMsg := fmt.Sprintf("Error writing to file: %s\n", err)
				printToBoth(writer, errorFileMsg)
			} else {
				printToBoth(writer, "Output file created: netfetch-cilium.txt\n")
			}
		} else {
			printToBoth(writer, "Output file not created.\n")
		}
	}

	score := CalculateScore(!missingPoliciesOrUncoveredPods, !userDeniedPolicyApplication, unprotectedPodsCount)
	scanResult.Score = score

	if printMessages {
		if policyChangesMade {
			fmt.Println("\nChanges were made during this scan. It's recommended to re-run the scan for an updated score.")
		}

		if missingPoliciesOrUncoveredPods {
			if userDeniedPolicyApplication {
				printToBoth(writer, "\nFor the following namespaces, you should assess the need of implementing network policies:\n")
				for _, ns := range deniedNamespaces {
					fmt.Println(" -", ns)
				}
				printToBoth(writer, "\nConsider either an implicit default deny all network policy or a policy that targets the pods not selected by a cilium network policy. Check the Kubernetes documentation for more information on cilium network policies: https://docs.cilium.io/en/latest/security/policy/\n")
			} else {
				printToBoth(writer, "\nNetfetch scan completed!\n")
			}
		} else {
			printToBoth(writer, "\nNo cilium network policies missing. You are good to go!\n")
		}
	}

	if printScore {
		// Print the final score
		fmt.Printf("\nYour Netfetch security score is: %d/42\n", score)
	}

	hasStartedCiliumScan = false
	return scanResult, nil
}

// convertEndpointToSelector converts the endpointSelector from a CiliumNetworkPolicy to a label selector string.
func convertEndpointToSelector(endpointSelector map[string]interface{}) (string, error) {
	matchLabels, found := endpointSelector["matchLabels"].(map[string]interface{})
	if !found || len(matchLabels) == 0 {
		return "", nil
	}

	var selectorParts []string
	for key, val := range matchLabels {
		if strVal, ok := val.(string); ok {
			selectorParts = append(selectorParts, fmt.Sprintf("%s=%s", key, strVal))
		} else {
			return "", fmt.Errorf("value for %s in matchLabels is not a string", key)
		}
	}

	return strings.Join(selectorParts, ","), nil
}

// createAndApplyDefaultDenyCiliumPolicy creates and applies a default deny all network policy for Cilium in the specified namespace.
func createAndApplyDefaultDenyCiliumPolicy(namespace string, dynamicClient dynamic.Interface) error {
	// Construct the policy name dynamically
	policyName := namespace + "-cilium-default-deny-all"

	// Define the policy
	denyAllPolicy := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cilium.io/v2",
			"kind":       "CiliumNetworkPolicy",
			"metadata": map[string]interface{}{
				"name":      policyName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"endpointSelector": map[string]interface{}{
					"matchLabels": map[string]string{},
				},
				"ingress": []interface{}{},
				"egress":  []interface{}{},
			},
		},
	}

	// Set the GVR
	ciliumNPResource := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v2",
		Resource: "ciliumnetworkpolicies",
	}

	// Create the policy
	_, err := dynamicClient.Resource(ciliumNPResource).Namespace(namespace).Create(context.TODO(), denyAllPolicy, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create default deny all CiliumNetworkPolicy: %v", err)
	}

	fmt.Printf("Applied default deny all CiliumNetworkPolicy to namespace %s\n", namespace)
	return nil
}

// hasDefaultDenyAllCiliumPolicy checks if the list of CiliumNetworkPolicies includes a default deny all policy
func hasDefaultDenyAllCiliumPolicy(policies []*unstructured.Unstructured) bool {
	for _, policy := range policies {
		if isDefaultDenyAllCiliumPolicy(*policy) {
			return true
		}
	}
	return false
}

// isDefaultDenyAllCiliumPolicy checks if a single Cilium policy is a default deny-all policy
func isDefaultDenyAllCiliumPolicy(policyUnstructured unstructured.Unstructured) bool {
	spec, found := policyUnstructured.UnstructuredContent()["spec"].(map[string]interface{})
	if !found {
		return false
	}

	ingress, ingressFound := spec["ingress"]
	egress, egressFound := spec["egress"]
	return (!ingressFound || len(ingress.([]interface{})) == 0) && (!egressFound || len(egress.([]interface{})) == 0)
}
