package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
func gatherVisualizationData(namespace string) (*VisualizationData, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, err
	}

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

// HandleVisualizationRequest handles the HTTP request for serving visualization data.
func HandleVisualizationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	namespace := r.URL.Query().Get("namespace")

	vizData, err := gatherVisualizationData(namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(vizData); err != nil {
		http.Error(w, "Failed to encode visualization data", http.StatusInternalServerError)
	}
}

// gatherNamespacesWithPolicies returns a list of all namespaces that contain network policies.
func GatherNamespacesWithPolicies() ([]string, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, err
	}

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
func GatherClusterVisualizationData() ([]VisualizationData, error) {
	namespacesWithPolicies, err := GatherNamespacesWithPolicies()
	if err != nil {
		return nil, err
	}

	// Slice to hold the visualization data for the entire cluster
	var clusterVizData []VisualizationData

	for _, namespace := range namespacesWithPolicies {
		vizData, err := gatherVisualizationData(namespace)
		if err != nil {
			log.Printf("Error gathering visualization data for namespace %s: %v\n", namespace, err)
			continue
		}
		clusterVizData = append(clusterVizData, *vizData)
	}

	return clusterVizData, nil
}

// structs to represent policy preview
type Port struct {
	Protocol string `yaml:"protocol,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	EndPort  int    `yaml:"endPort,omitempty"`
}

type IPBlock struct {
	CIDR string `yaml:"cidr,omitempty"`
}

type PolicyRule struct {
	From []struct {
		NamespaceSelector map[string]map[string]string `yaml:"namespaceSelector,omitempty"`
		PodSelector       map[string]map[string]string `yaml:"podSelector"`
	} `yaml:"from,omitempty"`
	To []struct {
		IPBlock *IPBlock `yaml:"ipBlock,omitempty"`
	} `yaml:"to,omitempty"`
	Ports []Port `yaml:"ports,omitempty"`
}

// struct for the NetworkPolicy
type NetworkPolicyPreview struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		PodSelector map[string]map[string]string `yaml:"podSelector"`
		Ingress     []PolicyRule                 `yaml:"ingress"`
		Egress      []PolicyRule                 `yaml:"egress"`
		PolicyTypes []string                     `yaml:"policyTypes"`
	} `yaml:"spec"`
}

// HandlePolicyYAMLRequest handles the HTTP request for serving the YAML of a network policy.
func HandlePolicyYAMLRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the policy name and namespace from query parameters
	policyName := r.URL.Query().Get("name")
	namespace := r.URL.Query().Get("namespace")
	if policyName == "" || namespace == "" {
		http.Error(w, "Policy name or namespace not provided", http.StatusBadRequest)
		return
	}

	// Retrieve the network policy YAML
	yaml, err := getNetworkPolicyYAML(namespace, policyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write([]byte(yaml))
}

// getNetworkPolicyYAML retrieves the YAML representation of a network policy.
func getNetworkPolicyYAML(namespace, policyName string) (string, error) {
	clientset, err := GetClientset()
	if err != nil {
		return "", err
	}

	// Get the specified network policy
	policy, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(context.TODO(), policyName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Convert the network policy into an unstructured object to easily extract fields
	unstructuredPolicy, err := runtime.DefaultUnstructuredConverter.ToUnstructured(policy)
	if err != nil {
		return "", err
	}

	// Extract the necessary fields from structs
	policyPreview := NetworkPolicyPreview{
		APIVersion: "networking.k8s.io/v1",
		Kind:       "NetworkPolicy",
	}

	// Metadata extraction
	unstructuredMeta, ok := unstructuredPolicy["metadata"].(map[string]interface{})
	if ok {
		policyPreview.Metadata.Name = unstructuredMeta["name"].(string)
		policyPreview.Metadata.Namespace = unstructuredMeta["namespace"].(string)
	}

	// Spec extraction
	spec, ok := unstructuredPolicy["spec"].(map[string]interface{})
	if ok {
		if podSelector, found := spec["podSelector"]; found {
			podSelectorMap, ok := podSelector.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("podSelector is not in expected format")
			}
			// Convert map[string]interface{} to map[string]map[string]string
			podSelectorStrMap := make(map[string]map[string]string)
			for k, v := range podSelectorMap {
				vStrMap, ok := v.(map[string]string)
				if !ok {
					return "", fmt.Errorf("podSelector value is not in expected format")
				}
				podSelectorStrMap[k] = vStrMap
			}
			policyPreview.Spec.PodSelector = podSelectorStrMap
		}
		var ingressPolicyRules []PolicyRule
		var egressPolicyRules []PolicyRule
		var err error

		if ingressRules, found := spec["ingress"]; found {
			ingressPolicyRules, err = extractPolicyRules(ingressRules)
			if err != nil {
				return "", err
			}
			policyPreview.Spec.Ingress = ingressPolicyRules
		}
		if egressRules, found := spec["egress"]; found {
			egressPolicyRules, err = extractPolicyRules(egressRules)
			if err != nil {
				return "", err
			}
			policyPreview.Spec.Egress = egressPolicyRules
		}
	}

	// Convert the network policy to YAML
	yamlBytes, err := yaml.Marshal(policyPreview)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

// extractPolicyRules extracts the rules from the ingress/egress section of the network policy
func extractPolicyRules(rules interface{}) ([]PolicyRule, error) {
	var policyRules []PolicyRule
	rulesSlice, ok := rules.([]interface{})
	if !ok {
		return nil, fmt.Errorf("rules are not in expected format")
	}

	for _, rule := range rulesSlice {
		ruleMap, ok := rule.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("rule is not in expected format")
		}

		policyRule := PolicyRule{}

		// Extract 'from' rules
		if fromRules, present := ruleMap["from"]; present {
			fromRulesSlice, ok := fromRules.([]interface{})
			if !ok {
				return nil, fmt.Errorf("'from' field is not in expected format")
			}
			for _, fromRule := range fromRulesSlice {
				fromMap, ok := fromRule.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("from rule is not in expected format")
				}
				var ruleFrom struct {
					NamespaceSelector map[string]map[string]string `yaml:"namespaceSelector,omitempty"`
					PodSelector       map[string]map[string]string `yaml:"podSelector"`
				}
				if nsSelector, nsFound := fromMap["namespaceSelector"]; nsFound {
					nsSelectorMap, ok := nsSelector.(map[string]map[string]string)
					if !ok {
						return nil, fmt.Errorf("namespaceSelector is not in expected format")
					}
					ruleFrom.NamespaceSelector = nsSelectorMap
				}
				if podSelector, podFound := fromMap["podSelector"]; podFound {
					podSelectorMap, ok := podSelector.(map[string]map[string]string)
					if !ok {
						return nil, fmt.Errorf("podSelector is not in expected format")
					}
					ruleFrom.PodSelector = podSelectorMap
				}
				policyRule.From = append(policyRule.From, ruleFrom)
			}
		}

		// Extract 'to' rules
		if toRules, present := ruleMap["to"]; present {
			toRulesSlice, ok := toRules.([]interface{})
			if !ok {
				return nil, fmt.Errorf("'to' field is not in expected format")
			}
			for _, toRule := range toRulesSlice {
				toMap, ok := toRule.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("to rule is not in expected format")
				}
				var ruleTo struct {
					IPBlock *IPBlock `yaml:"ipBlock,omitempty"`
				}
				if ipBlock, ipBlockFound := toMap["ipBlock"]; ipBlockFound {
					ipBlockMap, ok := ipBlock.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("ipBlock is not in expected format")
					}
					ruleTo.IPBlock = &IPBlock{
						CIDR: ipBlockMap["cidr"].(string),
					}
				}
				policyRule.To = append(policyRule.To, ruleTo)
			}
		}

		// Extract 'ports' rules
		if ports, present := ruleMap["ports"]; present {
			portsSlice, ok := ports.([]interface{})
			if !ok {
				return nil, fmt.Errorf("'ports' field is not in expected format")
			}
			for _, port := range portsSlice {
				portMap, ok := port.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("port is not in expected format")
				}
				var portStruct Port
				if protocol, ok := portMap["protocol"]; ok {
					portStruct.Protocol = protocol.(string)
				}
				if port, ok := portMap["port"]; ok {
					portStruct.Port = int(port.(float64)) // assuming port is provided as a float64
				}
				if endPort, ok := portMap["endPort"]; ok {
					portStruct.EndPort = int(endPort.(float64)) // assuming endPort is provided as a float64
				}
				policyRule.Ports = append(policyRule.Ports, portStruct)
			}
		}

		policyRules = append(policyRules, policyRule)
	}
	return policyRules, nil
}
