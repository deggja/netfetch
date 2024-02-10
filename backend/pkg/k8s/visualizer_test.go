package k8s

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes"
)

func TestGatherVisualizationData(t *testing.T) {
	var clientset kubernetes.Clientset
	clientset = fake.NewSimpleClientset()

	// Creating a pod object and a networkPolicy object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod", 
			Namespace: "test-namespace",
		},
	}
	networkPolicy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-policy", 
			Namespace: "test-namespace",
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
			},
		},
	}
	_, err := clientset.CoreV1().Pods("test-namespace").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create pod: %v", err)
	}
	_, erro := clientset.NetworkingV1().NetworkPolicies("test-namespace").Create(context.TODO(), networkPolicy, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create network policy: %v", erro)
	}

	visualizationData, err := gatherVisualizationData(&clientset, "test-namespace")
	if err != nil {
		t.Fatalf("Error occurred while gathering visualization data: %v", err)
	}

	expected := &VisualizationData{
		Policies: []PolicyVisualization{
			{
				Name:       "test-policy",
				Namespace:  "test-namespace",
				TargetPods: []string{"test-pod"},
			},
		},
	}
	if !reflect.DeepEqual(visualizationData, expected) {
		t.Errorf("Visualization data mismatch. Expected: %v, Got: %v", expected, visualizationData)
	}
}