package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestMatchesLabels(t *testing.T) {
	tests := []struct {
		name          string
		podLabels     map[string]string
		policyLabels  map[string]interface{}
		expectedMatch bool
	}{
		{
			name: "Matching labels",
			podLabels: map[string]string{
				"app": "test",
			},
			// Equivalent to MatchLabels in a NetworkPolicy object
			policyLabels: map[string]interface{}{
				"app": "test",
			},
			expectedMatch: true,
		},
		{
			name: "Matching labels",
			podLabels: map[string]string{
				"app": "test",
			},
			// Equivalent to MatchLabels in a NetworkPolicy object
			policyLabels: map[string]interface{}{
				"app": "test-policy",
			},
			expectedMatch: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := MatchesLabels(test.podLabels, test.policyLabels)
			// if match != test.expectedMatch {
			// 	t.Errorf("Expected match: %v, got: %v", test.expectedMatch, match)
			// }
			assert.Equal(t, test.expectedMatch, match, "they should be equal")
		})
	}
}

func TestConvertEndpointToSelector(t *testing.T) {
	tests := []struct {
		name              string
		endpointSelector  map[string]interface{}
		expectedSelector  string
		expectedError     error
	}{
		{
			name: "Non-Empty Selector",
			endpointSelector: map[string]interface{}{
				"matchLabels": map[string]interface{}{"app": "test"},
			},
			expectedSelector: "app=test",
			expectedError: nil,
		},
		{
			name: "Empty Selector",
			endpointSelector: map[string]interface{}{
				"matchLabels": map[string]string{},
			},
			expectedSelector: "",
			expectedError:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector, err := ConvertEndpointToSelector(test.endpointSelector)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assert.Equal(t, test.expectedSelector, selector)
		})
	}
}

func TestIsDefaultDenyAllCiliumClusterwidePolicy(t *testing.T) {
	tests := []struct {
		name               string
		policyUnstructured unstructured.Unstructured
		expectedDenyAll    bool
		expectedCluster    bool
	}{
		{
			name: "Default Deny All Policy",
			policyUnstructured: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"ingress": []interface{}{},
						"egress":  []interface{}{},
						"endpointSelector": map[string]interface{}{
							"matchLabels": map[string]interface{}{},
						},
					},
				},
			},
			expectedDenyAll: true,
			expectedCluster: true,
		},
		{
			name: "Non-Default Policy",
			policyUnstructured: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"ingress": []interface{}{
							map[string]interface{}{
								"port": 80,
							}},
						"egress":  []interface{}{},
						"endpointSelector": map[string]interface{}{
							"matchLabels": map[string]interface{}{"app": "test"},
						},
					},
				},
			},
			expectedDenyAll: false,
			expectedCluster: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			denyAll, cluster := IsDefaultDenyAllCiliumClusterwidePolicy(test.policyUnstructured)

			if denyAll != test.expectedDenyAll {
				t.Errorf("Expected Default Deny All: %v, got: %v", test.expectedDenyAll, denyAll)
			}

			if cluster != test.expectedCluster {
				t.Errorf("Expected Clusterwide: %v, got: %v", test.expectedCluster, cluster)
			}
		})
	}
}

func TestIsEmptyOrOnlyContainsEmptyObjects(t *testing.T) {
	tests := []struct {
		name          string
		slice         []interface{}
		expectedEmpty bool
	}{
		{
			name:          "Empty Slice",
			slice:         []interface{}{},
			expectedEmpty: true,
		},
		{
			name:          "Non-Empty Slice",
			slice:         []interface{}{map[string]interface{}{"key": "value"}},
			expectedEmpty: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isEmpty := IsEmptyOrOnlyContainsEmptyObjects(test.slice)

			if isEmpty != test.expectedEmpty {
				t.Errorf("Expected empty: %v, got: %v", test.expectedEmpty, isEmpty)
			}
		})
	}
}

func TestIsEndpointSelectorEmpty(t *testing.T) {
	tests := []struct {
		name string
		selector map[string]interface{}
		expectedIsEmpty bool
	}{
		{
			name: "Empty Selector",
			selector: map[string]interface{}{},
			expectedIsEmpty: true,
		},
		{
			name: "Non-Empty Selector",
			selector: map[string]interface{}{
				"matchLabels": map[string]interface{}{"app": "test"},
			},
			expectedIsEmpty: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isEmpty := isEndpointSelectorEmpty(test.selector)

			if isEmpty != test.expectedIsEmpty {
				t.Errorf("Expected is empty: %v, got: %v", test.expectedIsEmpty, isEmpty)
			}
		})
	}
}

func TestIsDefaultDenyAllCiliumPolicy(t *testing.T) {
	tests := []struct {
		name               string
		policyUnstructured unstructured.Unstructured
		expectedDenyAll    bool
	}{
		{
			name: "Default Deny All Policy",
			policyUnstructured: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"ingress": []interface{}{},
						"egress":  []interface{}{},
					},
				},
			},
			expectedDenyAll: true,
		},
		{
			name: "Non-Default Policy",
			policyUnstructured: unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"ingress": []interface{}{map[string]interface{}{"port": 80}},
						"egress":  []interface{}{},
					},
				},
			},
			expectedDenyAll: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			denyAll := IsDefaultDenyAllCiliumPolicy(test.policyUnstructured)

			if denyAll != test.expectedDenyAll {
				t.Errorf("Expected Default Deny All: %v, got: %v", test.expectedDenyAll, denyAll)
			}
		})
	}
}