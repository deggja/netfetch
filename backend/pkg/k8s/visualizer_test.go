package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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

	var ns string = "test-namespace"

	// Creating a pod object and a networkPolicy object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Namespace: ns,
			Labels: map[string]string{
				"app": "test", // Label on pod to match netpol selector
			},
		},
	}
	_, err := clientset.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create pod: %v", err)
	}

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
	_, err1 := clientset.NetworkingV1().NetworkPolicies(ns).Create(context.TODO(), testPolicy, metav1.CreateOptions{})
	if err1 != nil {
		t.Fatalf("Failed to create test network policy: %v", err1)
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
	err2 := yaml.Unmarshal(expectedPolicyBytes, &expectedPolicy)
	if err2 != nil {
		t.Fatalf("error: %v", err2)
	}

	var actualPolicy Policy
	err3 := yaml.Unmarshal([]byte(YAML), &actualPolicy)
	if err3 != nil {
		t.Fatalf("error: %v", err3)
	}

	assert.Equal(t, expectedPolicy, actualPolicy, "they should be equal")
}

func TestGatherNamespacesWithPolicies(t *testing.T) {
	// Define namespaces both with and without network policies
	namespaces := []string{"test-namespace1", "test-namespace2", "test-namespace3", "test-namespace4"}
	namespacesWithPolicies := []string{"test-namespace3", "test-namespace4"}

	var clientset kubernetes.Interface = fake.NewSimpleClientset()

	for _, ns := range namespaces {
		// Create a namespace
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		}
		
		_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Failed to create namespace %s: %v", ns, err)
		}
	}

	for _, ns := range namespacesWithPolicies {
		// Create a pod in the namespace
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pod",
				Namespace: ns,
				Labels: map[string]string{
					"app": "test",
				},
			},
		}
		_, err := clientset.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Failed to create pod: %v", err)
		}

		// Create a networkPolicy in the namespace
		networkPolicy := &netv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-policy",
				Namespace: ns,
			},
			Spec: netv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "test"},
				},
			},
		}
		_, err1 := clientset.NetworkingV1().NetworkPolicies(ns).Create(context.TODO(), networkPolicy, metav1.CreateOptions{})
		if err1 != nil {
			t.Fatalf("Failed to create network policy in namespace %s: %v", ns, err1)
		}
	}

	gatheredNamespaces, err := GatherNamespacesWithPolicies(clientset)
	if err != nil {
		t.Fatalf("Error calling GatherNamespacesWithPolicies: %v", err)
	}

	assert.Equal(t, namespacesWithPolicies, gatheredNamespaces, "they should be equal")
}

func TestGatherClusterVisualizationData(t *testing.T) {
    namespacesWithPolicies := []string{"namespace1", "namespace2"}
    expectedVisualizationData := make([]VisualizationData, len(namespacesWithPolicies))

	var clientset kubernetes.Interface = fake.NewSimpleClientset()

    for i, ns := range namespacesWithPolicies {
		// Create a namespace
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		}
		_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Failed to create namespace %s: %v", ns, err)
		}

		// Create a pod in the namespace
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pod",
				Namespace: ns,
				Labels: map[string]string{
					"app": "test",
				},
			},
		}
		_, err1 := clientset.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err1 != nil {
			t.Fatalf("Failed to create pod: %v", err1)
		}

        // Create a network policy in the namespace
        networkPolicy := &netv1.NetworkPolicy{
            ObjectMeta: metav1.ObjectMeta{
				Name: "test-policy",
				Namespace: ns,
			},
			Spec: netv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "test"},
				},
			},
        }
        _, err2 := clientset.NetworkingV1().NetworkPolicies(ns).Create(context.TODO(), networkPolicy, metav1.CreateOptions{})
        if err2 != nil {
            t.Fatalf("Failed to create network policy in namespace %s: %v", ns, err2)
        }

        vizData := &VisualizationData{
            Policies: []PolicyVisualization{
                {
                    Name:       networkPolicy.Name,
                    Namespace:  ns,
					TargetPods: []string{"test-pod"}, 
                },
            },
        }
        expectedVisualizationData[i] = *vizData
    }

    gatheredVisualizationData, err := GatherClusterVisualizationData(clientset)
    if err != nil {
        t.Fatalf("Error calling GatherClusterVisualizationData: %v", err)
    }

    assert.Equal(t, expectedVisualizationData, gatheredVisualizationData)
}