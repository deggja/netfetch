package k8s

import (
	"context"
	"fmt"
	"regexp"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// FindNativeNetworkPolicyByName searches for a specific native network policy by name across all non-system namespaces.
func FindNativeNetworkPolicyByName(dynamicClient dynamic.Interface, clientset *kubernetes.Clientset, policyName string) (*unstructured.Unstructured, string, error) {
	gvr := schema.GroupVersionResource{
		Group:    "networking.k8s.io",
		Version:  "v1",
		Resource: "networkpolicies",
	}

	namespaces, err := GetAllNonSystemNamespaces(dynamicClient)
	if err != nil {
		return nil, "", fmt.Errorf("error getting namespaces: %v", err)
	}

	for _, namespace := range namespaces {
		policy, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), policyName, v1.GetOptions{})
		if err == nil {
			return policy, namespace, nil
		}
	}
	return nil, "", fmt.Errorf("network policy %s not found in any non-system namespace", policyName)
}

// FindCiliumNetworkPolicyByName searches for a specific Cilium network policy by name across all non-system namespaces.
func FindCiliumNetworkPolicyByName(dynamicClient dynamic.Interface, policyName string) (*unstructured.Unstructured, string, error) {
    gvr := schema.GroupVersionResource{
        Group:    "cilium.io",
        Version:  "v2",
        Resource: "ciliumnetworkpolicies",
    }

    namespaces, err := GetAllNonSystemNamespaces(dynamicClient)
    if err != nil {
        return nil, "", fmt.Errorf("error getting namespaces: %v", err)
    }

    for _, namespace := range namespaces {
        policy, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), policyName, v1.GetOptions{})
        if err == nil {
            return policy, namespace, nil
        }
    }
    return nil, "", fmt.Errorf("cilium network policy %s not found in any non-system namespace", policyName)
}

// FindCiliumClusterWideNetworkPolicyByName searches for a specific cluster wide Cilium network policy by name.
func FindCiliumClusterWideNetworkPolicyByName(dynamicClient dynamic.Interface, policyName string) (*unstructured.Unstructured, error) {
    gvr := schema.GroupVersionResource{
        Group:    "cilium.io",
        Version:  "v2",
        Resource: "ciliumclusterwidenetworkpolicies",
    }

    policy, err := dynamicClient.Resource(gvr).Get(context.TODO(), policyName, v1.GetOptions{})
    if err != nil {
        return nil, fmt.Errorf("cilium cluster wide network policy %s not found", policyName)
    }
    return policy, nil
}


// GetAllNonSystemNamespaces returns a list of all non-system namespaces using a dynamic client.
func GetAllNonSystemNamespaces(dynamicClient dynamic.Interface) ([]string, error) {
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}

	namespacesList, err := dynamicClient.Resource(gvr).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing namespaces: %v", err)
	}

	var namespaces []string
	for _, ns := range namespacesList.Items {
		if !IsSystemNamespace(ns.GetName()) {
			namespaces = append(namespaces, ns.GetName())
		}
	}
	return namespaces, nil
}

// ListPodsTargetedByNetworkPolicy lists all pods targeted by the given network policy in the specified namespace.
func ListPodsTargetedByNetworkPolicy(dynamicClient dynamic.Interface, policy *unstructured.Unstructured, namespace string) ([][]string, error) {
	// Retrieve the PodSelector (matchLabels)
	podSelector, found, err := unstructured.NestedMap(policy.Object, "spec", "podSelector", "matchLabels")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pod selector from network policy %s: %v", policy.GetName(), err)
	}

	// Check if the selector is empty
	selector := make(labels.Set)
	if found && len(podSelector) > 0 {
		for key, value := range podSelector {
			if strValue, ok := value.(string); ok {
				selector[key] = strValue
			} else {
				return nil, fmt.Errorf("invalid type for selector value %v in policy %s", value, policy.GetName())
			}
		}
	}

	// Fetch pods based on the selector
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{LabelSelector: selector.AsSelectorPreValidated().String()})
	if err != nil {
		return nil, fmt.Errorf("error listing pods in namespace %s: %v", namespace, err)
	}

	var targetedPods [][]string
    for _, pod := range pods.Items {
        podDetails := []string{namespace, pod.Name, pod.Status.PodIP}
        if pod.Status.PodIP == "" {
            podDetails[2] = "N/A"
        }
        targetedPods = append(targetedPods, podDetails)
    }

	return targetedPods, nil
}

// ListPodsTargetedByCiliumNetworkPolicy lists all pods targeted by the given Cilium network policy in the specified namespace.
func ListPodsTargetedByCiliumNetworkPolicy(dynamicClient dynamic.Interface, policy *unstructured.Unstructured, namespace string) ([][]string, error) {
    // Retrieve the PodSelector (matchLabels)
    podSelector, found, err := unstructured.NestedMap(policy.Object, "spec", "endpointSelector", "matchLabels")
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve pod selector from Cilium network policy %s: %v", policy.GetName(), err)
    }

    // Check if the selector is empty
    selector := make(labels.Set)
    if found && len(podSelector) > 0 {
        for key, value := range podSelector {
            if strValue, ok := value.(string); ok {
                selector[key] = strValue
            } else {
                return nil, fmt.Errorf("invalid type for selector value %v in policy %s", value, policy.GetName())
            }
        }
    }

    // Fetch pods based on the selector
    pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{LabelSelector: selector.AsSelectorPreValidated().String()})
    if err != nil {
        return nil, fmt.Errorf("error listing pods in namespace %s: %v", namespace, err)
    }

    var targetedPods [][]string
    for _, pod := range pods.Items {
        targetedPods = append(targetedPods, []string{namespace, pod.Name, pod.Status.PodIP})
    }

    return targetedPods, nil
}

// ListPodsTargetedByCiliumClusterWideNetworkPolicy lists all pods targeted by the given Cilium cluster wide network policy.
func ListPodsTargetedByCiliumClusterWideNetworkPolicy(clientset *kubernetes.Clientset, dynamicClient dynamic.Interface, policy *unstructured.Unstructured) ([][]string, error) {
    // Retrieve the PodSelector (matchLabels)
    podSelector, found, err := unstructured.NestedMap(policy.Object, "spec", "endpointSelector", "matchLabels")
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve pod selector from Cilium cluster wide network policy %s: %v", policy.GetName(), err)
    }

    // Regex for valid Kubernetes label keys
    validLabelKey := regexp.MustCompile(`^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$`)

    // Check if the selector is empty
    selector := labels.Set{}
    if found && len(podSelector) > 0 {
        for key, value := range podSelector {
            // Skip reserved labels
            if !validLabelKey.MatchString(key) {
                fmt.Printf("Skipping reserved label key %s in policy %s\n", key, policy.GetName())
                continue
            }
            if strValue, ok := value.(string); ok {
                selector[key] = strValue
            } else {
                return nil, fmt.Errorf("invalid type for selector value %v in policy %s", value, policy.GetName())
            }
        }
    }

    // Fetch pods based on the selector across all namespaces
    pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{
        LabelSelector: selector.AsSelector().String(),
    })
    if err != nil {
        return nil, fmt.Errorf("error listing pods for cluster wide policy: %v", err)
    }

    var targetedPods [][]string
    for _, pod := range pods.Items {
        podDetails := []string{pod.Namespace, pod.Name, pod.Status.PodIP}
        targetedPods = append(targetedPods, podDetails)
    }

    return targetedPods, nil
}