package k8s

import (
	"context"
	"log"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// VisualizationData represents the structure of network policy and pod data for visualization.
type VisualizationData struct {
	Policies []PolicyVisualization `json:"policies"`
}

// PolicyVisualization represents a network policy and the pods it affects for visualization purposes.
type PolicyVisualization struct {
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	TargetPods []string `json:"targetPods"`
}

// gatherVisualizationData retrieves network policies and associated pods for visualization.
func gatherVisualizationData(clientset kubernetes.Interface, namespace string) (*VisualizationData, error) {
	// Retrieve all network policies in the specified namespace
	policies, err := clientset.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	vizData := &VisualizationData{
		Policies: make([]PolicyVisualization, 0), // Initialize as empty slice
	}

	// Iterate over the retrieved policies to build the visualization data
	for _, policy := range policies.Items {
		// For each policy, find the pods that match its pod selector
		selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.PodSelector)
		if err != nil {
			log.Printf("Error parsing selector for policy %s: %v\n", policy.Name, err)
			continue
		}

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		if err != nil {
			log.Printf("Error listing pods for policy %s: %v\n", policy.Name, err)
			continue
		}

		podNames := make([]string, 0, len(pods.Items))
		for _, pod := range pods.Items {
			podNames = append(podNames, pod.Name)
		}

		vizData.Policies = append(vizData.Policies, PolicyVisualization{
			Name:       policy.Name,
			Namespace:  policy.Namespace,
			TargetPods: podNames,
		})
	}

	return vizData, nil
}

// gatherNamespacesWithPolicies returns a list of all namespaces that contain network policies.
func GatherNamespacesWithPolicies(clientset kubernetes.Interface) ([]string, error) {
	// Retrieve all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespacesWithPolicies []string

	// Check each namespace for network policies
	for _, ns := range namespaces.Items {
		policies, err := clientset.NetworkingV1().NetworkPolicies(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error listing policies in namespace %s: %v\n", ns.Name, err)
			continue
		}

		if len(policies.Items) > 0 {
			namespacesWithPolicies = append(namespacesWithPolicies, ns.Name)
		}
	}

	return namespacesWithPolicies, nil
}

// gatherClusterVisualizationData retrieves visualization data for all namespaces with network policies.
func GatherClusterVisualizationData(clientset kubernetes.Interface) ([]VisualizationData, error) {
	namespacesWithPolicies, err := GatherNamespacesWithPolicies(clientset)
	if err != nil {
		return nil, err
	}

	// Slice to hold the visualization data for the entire cluster
	var clusterVizData []VisualizationData

	for _, namespace := range namespacesWithPolicies {
		vizData, err := gatherVisualizationData(clientset, namespace)
		if err != nil {
			log.Printf("Error gathering visualization data for namespace %s: %v\n", namespace, err)
			continue
		}
		clusterVizData = append(clusterVizData, *vizData)
	}

	return clusterVizData, nil
}

// getNetworkPolicyYAML retrieves the YAML representation of a network policy, excluding annotations.
func getNetworkPolicyYAML(clientset kubernetes.Interface, namespace string, policyName string) (string, error) {
	// Get the specified network policy
	networkPolicy, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(context.TODO(), policyName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Convert the network policy to a map[string]interface{} for manipulation
	networkPolicyMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(networkPolicy)
	if err != nil {
		return "", err
	}

	// Remove the fields that are not needed
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "annotations")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "generateName")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "selfLink")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "uid")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "generation")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "creationTimestamp")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "deletionTimestamp")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "deletionGracePeriodSeconds")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "ownerReferences")
	unstructured.RemoveNestedField(networkPolicyMap, "metadata", "managedFields")

	// Convert the cleaned map back to YAML
	yamlBytes, err := yaml.Marshal(networkPolicyMap)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}
