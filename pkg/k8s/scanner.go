package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ScanNetworkPolicies scans all non-system namespaces for network policies
func ScanNetworkPolicies() {
	// Use the default kubeconfig path if running outside the cluster
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Use the current context in kubeconfig
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

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing namespaces: %s\n", err)
		return
	}

	for _, ns := range namespaces.Items {
		if isSystemNamespace(ns.Name) {
			continue
		}

		policies, err := clientset.NetworkingV1().NetworkPolicies(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("\nError listing network policies in namespace %s: %s\n", ns.Name, err)
			continue
		}

		coveredPods := make(map[string]bool)
		hasPolicies := false

		for _, policy := range policies.Items {
			hasPolicies = true
			fmt.Printf("Network policy found: %s in namespace %s\n", policy.Name, ns.Name)

			// Get the pods targeted by this policy
			selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.PodSelector)
			if err != nil {
				fmt.Printf("Error parsing selector for policy %s: %s\n", policy.Name, err)
				continue
			}

			pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{
				LabelSelector: selector.String(),
			})
			if err != nil {
				fmt.Printf("Error listing pods for policy %s: %s\n", policy.Name, err)
				continue
			}

			for _, pod := range pods.Items {
				coveredPods[pod.Name] = true
				fmt.Printf("Pod covered by policy %s: Name: %s, IP: %s\n", policy.Name, pod.Name, pod.Status.PodIP)
			}
		}

		if !hasPolicies || !hasDefaultDenyAllPolicy(policies.Items) {
			allPods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				fmt.Printf("Error listing all pods in namespace %s: %s\n", ns.Name, err)
				continue
			}

			unprotectedPods := false
			for _, pod := range allPods.Items {
				if !coveredPods[pod.Name] {
					unprotectedPods = true
					fmt.Printf("Unprotected Pod: Name: %s, IP: %s\n", pod.Name, pod.Status.PodIP)
				}
			}

			if unprotectedPods {
				// Ask the user whether to add a network policy
				confirm := false
				prompt := &survey.Confirm{
					Message: fmt.Sprintf("Do you want to add a default deny all network policy to the namespace %s?", ns.Name),
				}
				survey.AskOne(prompt, &confirm, nil)

				if confirm {
					err := createAndApplyDefaultDenyPolicy(clientset, ns.Name)
					if err != nil {
						fmt.Printf("\nFailed to apply default deny policy in namespace %s: %s\n", ns.Name, err)
					} else {
						fmt.Printf("\nApplied default deny policy in namespace %s\n", ns.Name)
					}
				}
			}
		}
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
