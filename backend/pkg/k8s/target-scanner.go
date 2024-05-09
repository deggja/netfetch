package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// GetAllNonSystemNamespaces returns a list of all non-system namespaces using a dynamic client.
func GetAllNonSystemNamespaces(dynamicClient dynamic.Interface) ([]string, error) {
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}

	namespacesList, err := dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
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
func ListPodsTargetedByNetworkPolicy(cynamicClient dynamic.Interface, policy *unstructured.Unstructured, namespace string) ([]string, error) {
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
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.AsSelectorPreValidated().String()})
	if err != nil {
		return nil, fmt.Errorf("error listing pods in namespace %s: %v", namespace, err)
	}

	var targetedPods []string
	for _, pod := range pods.Items {
		targetedPods = append(targetedPods, pod.Name)
	}

	return targetedPods, nil
}
