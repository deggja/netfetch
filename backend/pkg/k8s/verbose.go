package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

// LogPolicyPodCoverage logs which pods are covered by a given policy.
func LogPolicyPodCoverage(clientset *kubernetes.Clientset, policy *unstructured.Unstructured, namespace string) error {
	policyName := policy.GetName()

	spec, found := policy.UnstructuredContent()["spec"].(map[string]interface{})
	if !found {
		return fmt.Errorf("policy %s does not have a spec", policyName)
	}

	endpointSelector, found := spec["endpointSelector"].(map[string]interface{})
	if !found {
		return fmt.Errorf("policy %s does not have an endpointSelector", policyName)
	}

	labelSelector, err := ConvertEndpointToSelector(endpointSelector)
	if err != nil {
		return fmt.Errorf("error converting endpoint selector for policy %s: %s", policyName, err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return fmt.Errorf("error listing pods for policy %s: %s", policyName, err)
	}

	fmt.Printf("Policy %s covers the following pods in namespace %s:\n", policyName, namespace)
	for _, pod := range pods.Items {
		fmt.Printf("- Pod name: %s\n", pod.Name)
	}

	return nil
}
