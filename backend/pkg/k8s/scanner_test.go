package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestHasDefaultDenyAllPolicy(t *testing.T) {
    // Test case 1: Default deny all policy exists
    policyWithDefaultDeny := netv1.NetworkPolicy{
        Spec: netv1.NetworkPolicySpec{},
    }
    if !hasDefaultDenyAllPolicy([]netv1.NetworkPolicy{policyWithDefaultDeny}) {
        t.Errorf("Expected to identify default deny all policy, but it was not detected")
    }

    // Test case 2: No default deny all policy
    policyWithoutDefaultDeny := netv1.NetworkPolicy{
        Spec: netv1.NetworkPolicySpec{
            Ingress: []netv1.NetworkPolicyIngressRule{
                {},
            },
            Egress: []netv1.NetworkPolicyEgressRule{
                {},
            },
        },
    }
    if hasDefaultDenyAllPolicy([]netv1.NetworkPolicy{policyWithoutDefaultDeny}) {
        t.Errorf("Expected not to identify default deny all policy, but it was detected")
    }
}

func TestIsDefaultDenyAllPolicy(t *testing.T) {
	// Test case 1: Default deny all policy
	defaultDenyPolicy := netv1.NetworkPolicy{
		Spec: netv1.NetworkPolicySpec{},
	}
	if !isDefaultDenyAllPolicy(defaultDenyPolicy) {
		t.Fatalf("Expected policy to be default deny all, but it was not")
	}

	// Test case 2: Non-default deny all policy
	nonDefaultDenyPolicy := netv1.NetworkPolicy{
		Spec: netv1.NetworkPolicySpec{
			Ingress: []netv1.NetworkPolicyIngressRule{
				{},
			},
			Egress: []netv1.NetworkPolicyEgressRule{
				{},
			},
		},
	}
	if isDefaultDenyAllPolicy(nonDefaultDenyPolicy) {
		t.Fatalf("Expected policy not to be default deny all, but it was")
	}
}

func TestIsSystemNamespace(t *testing.T) {
	// Test case 1: System namespace
	systemNamespaces := []string{"kube-system", "tigera-operator", "kube-public", "kube-node-lease", "gatekeeper-system", "calico-system"}
	for _, ns := range systemNamespaces {
		if !IsSystemNamespace(ns) {
			t.Fatalf("Expected namespace %s to be a system namespace, but it was not", ns)
		}
	}

	// Test case 2: Non-system namespace
	nonSystemNamespace := "test-namespace"
	if IsSystemNamespace(nonSystemNamespace) {
		t.Fatalf("Expected namespace %s not to be a system namespace, but it was", nonSystemNamespace)
	}
}

func TestCalculateScore(t *testing.T) {
	// Test case 1: All conditions met, maximum score expected
	score1 := CalculateScore(true, true, 0)
	if score1 != 42 {
		t.Fatalf("Expected score to be 42, got %d", score1)
	}

	// Test case 2: No policies, no deny all, 5 unprotected pods
	score2 := CalculateScore(false, false, 5)
	if score2 != 17 {
		t.Fatalf("Expected score to be 17, got %d", score2)
	}

	// Test case 3: No policies, no deny all, no unprotected pods
	score3 := CalculateScore(false, false, 0)
	if score3 != 22 {
		t.Fatalf("Expected score to be 1, got %d", score3)
	}

	// Test case 4: Policies exist, deny all exists, no unprotected pods
	score4 := CalculateScore(true, true, 0)
	if score4 != 42 {
		t.Fatalf("Expected score to be 42, got %d", score4)
	}
}

func TestGetPodInfo(t *testing.T) {
    var clientset kubernetes.Interface = fake.NewSimpleClientset()

    podInfo := PodInfo{
        Name: "test-pod",
        Namespace: "test-namespace",
        Labels: map[string]string{
            "app": "test", // Label on pod to match netpol selector
        },
        Ports: []corev1.ContainerPort{
            {
                ContainerPort: 80,
            },
        },
    }
	
    pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podInfo.Name,
			Namespace: podInfo.Namespace,
			Labels:    podInfo.Labels,
		},
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{
                {
                    Name:  "nginx",
                    Image: "nginx",
                    Ports: podInfo.Ports,
                },
            },
        },
	}
    _, err := clientset.CoreV1().Pods("test-namespace").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create pod: %v", err)
	}

	var expectedPodInfo []PodInfo
	expectedPodInfo = append(expectedPodInfo, podInfo)

    actualPodInfo, err := GetPodInfo(clientset, podInfo.Namespace)
	if err != nil {
		t.Fatalf("Failed to get actual podInfo: %v", err)
	}

    assert.Equal(t, expectedPodInfo, actualPodInfo, "they should be equal")
}

func TestYAMLToNetworkPolicy(t *testing.T) {
	YAMLString := `
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-policy
  namespace: test-namespace
spec:
  podSelector:
    matchLabels:
      app: test
`
	expectedNetworkPolicy := &netv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "NetworkPolicy",
		},
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

	actualNetworkPolicy, err := YAMLToNetworkPolicy(YAMLString)
	if err != nil {
		t.Fatalf("Failed to convert YAML string to a NetworkPolicy object: %v", err)
	}

	assert.Equal(t, expectedNetworkPolicy, actualNetworkPolicy, "they should be equal")
}