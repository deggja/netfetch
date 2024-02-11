package k8s

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type PodSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type Policy struct {
	Metadata    Metadata    `yaml:"metadata"`
	Spec        struct {
		PodSelector PodSelector `yaml:"podSelector"`
	} `yaml:"spec"`
}

func TestGatherVisualizationData(t *testing.T) {
	var clientset kubernetes.Interface = fake.NewSimpleClientset()

	// Creating a pod object and a networkPolicy object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
			Labels: map[string]string{
				"app": "test", // Label on pod to match netpol selector
			},
		},
	}
	networkPolicy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-policy",
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
	_, err = clientset.NetworkingV1().NetworkPolicies("test-namespace").Create(context.TODO(), networkPolicy, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create network policy: %v", err)
	}

	visualizationData, err := gatherVisualizationData(clientset, "test-namespace")
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

	assert.Equal(t, expected, visualizationData, "they should be equal")
}

func TestGetNetworkPolicyYAML(t *testing.T) {
	var clientset kubernetes.Interface = fake.NewSimpleClientset()

	testPolicy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-policy",
			Namespace: "test-namespace",
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
			},
		},
	}

	_, err := clientset.NetworkingV1().NetworkPolicies("test-namespace").Create(context.TODO(), testPolicy, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test network policy: %v", err)
	}

	// Call the getNetworkPolicyYAML function with the fake clientset
	YAML, err := getNetworkPolicyYAML(clientset, "test-namespace", "test-policy")
	if err != nil {
		t.Fatalf("Failed to get network policy YAML: %v", err)
	}

	// Policy YAML string
	// expectedYAML := `
	// metadata:
	//   name: test-policy
	//   namespace: test-namespace
	// spec:
	//   podSelector:
	//     matchLabels:
	// 	  app: test
	// `

	// Below the testPolicy which is expected is converted into Policy struct (map to bytes array to struct) and
	// compared with actualPolicy which is also converted to Policy struct from string.
	expectedNetworkPolicyMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(testPolicy)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	expectedPolicyBytes, err := yaml.Marshal(expectedNetworkPolicyMap)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var expectedPolicy Policy
	err1 := yaml.Unmarshal(expectedPolicyBytes, &expectedPolicy)
	if err1 != nil {
		t.Fatalf("error: %v", err1)
	}

	var actualPolicy Policy
	err2 := yaml.Unmarshal([]byte(YAML), &actualPolicy)
	if err2 != nil {
		t.Fatalf("error: %v", err2)
	}

	assert.Equal(t, expectedPolicy, actualPolicy, "they should be equal")
}