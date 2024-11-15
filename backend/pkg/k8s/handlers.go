package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// INTERACTIVE DASHBOARD LOGIC

func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// HandleScanRequest handles the HTTP request for scanning network policies
func HandleScanRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        namespace := r.URL.Query().Get("namespace")

        // Perform the scan
        result, err := ScanNetworkPolicies(namespace, false, true, false, false, false, kubeconfigPath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Respond with JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}

// HandleNamespaceListRequest lists all non-system Kubernetes namespaces
func HandleNamespaceListRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, "Failed to create Kubernetes client: "+err.Error(), http.StatusInternalServerError)
            return
        }

        namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
        if err != nil {
            // Handle forbidden access error specifically
            if statusErr, isStatus := err.(*k8serrors.StatusError); isStatus {
                if statusErr.Status().Code == http.StatusForbidden {
                    http.Error(w, "Access forbidden: "+err.Error(), http.StatusForbidden)
                    return
                }
            }
            http.Error(w, "Failed to list namespaces: "+err.Error(), http.StatusInternalServerError)
            return
        }

        var namespaceList []string
        for _, ns := range namespaces.Items {
            if !IsSystemNamespace(ns.Name) {
                namespaceList = append(namespaceList, ns.Name)
            }
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string][]string{"namespaces": namespaceList})
    }
}

func HandleAddPolicyRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Define a struct to parse the incoming request
        type request struct {
            Namespace string `json:"namespace"`
        }

        // Parse the incoming JSON request
        var req request
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Apply the default deny policy
        err := createAndApplyDefaultDenyPolicy(req.Namespace, kubeconfigPath)
        if err != nil {
            http.Error(w, "Failed to apply default deny policy: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Respond with success message
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Implicit default deny all network policy successfully added to namespace " + req.Namespace,
        })

        // Re-scan the namespace
        scanResult, err := ScanNetworkPolicies(req.Namespace, false, true, false, false, false, kubeconfigPath)
        if err != nil {
            http.Error(w, "Error re-scanning after applying policy: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Respond with updated scan results
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(scanResult)
    }
}

// HandleNamespacesWithPoliciesRequest handles the HTTP request for serving a list of namespaces with network policies.
func HandleNamespacesWithPoliciesRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, "You are not connected to a Kubernetes cluster. Please connect to a cluster and re-run the command: "+err.Error(), http.StatusInternalServerError)
            return
        }

        namespaces, err := GatherNamespacesWithPolicies(clientset)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        setNoCacheHeaders(w)
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(struct {
            Namespaces []string `json:"namespaces"`
        }{Namespaces: namespaces}); err != nil {
            http.Error(w, "Failed to encode namespaces data", http.StatusInternalServerError)
        }
    }
}

// HandleNamespacePoliciesRequest handles the HTTP request for serving a list of network policies in a namespace.
func HandleNamespacePoliciesRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        // Extract the namespace parameter from the query string
        namespace := r.URL.Query().Get("namespace")
        if namespace == "" {
            http.Error(w, "Namespace parameter is required", http.StatusBadRequest)
            return
        }

        // Obtain the Kubernetes clientset
        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to create Kubernetes client: %v", err), http.StatusInternalServerError)
            return
        }

        // Fetch network policies from the specified namespace
        policies, err := clientset.NetworkingV1().NetworkPolicies(namespace).List(context.Background(), metav1.ListOptions{})
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to get network policies: %v", err), http.StatusInternalServerError)
            return
        }

        setNoCacheHeaders(w)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(policies)
    }
}

// HandleClusterVisualizationRequest handles the HTTP request for serving cluster-wide visualization data.
func HandleClusterVisualizationRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Call the function to gather cluster-wide visualization data
        clusterVizData, err := GatherClusterVisualizationData(clientset)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        setNoCacheHeaders(w)
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(clusterVizData); err != nil {
            http.Error(w, "Failed to encode cluster visualization data", http.StatusInternalServerError)
        }
    }
}

// HandlePodInfoRequest handles the HTTP request for serving pod information.
func HandlePodInfoRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        // Extract the namespace parameter from the query string
        namespace := r.URL.Query().Get("namespace")
        if namespace == "" {
            http.Error(w, "Namespace parameter is required", http.StatusBadRequest)
            return
        }

        // Obtain the Kubernetes clientset
        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to create Kubernetes client: %v", err), http.StatusInternalServerError)
            return
        }

        // Fetch pod information from the specified namespace
        podInfo, err := GetPodInfo(clientset, namespace)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to get pod information: %v", err), http.StatusInternalServerError)
            return
        }

        setNoCacheHeaders(w)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(podInfo)
    }
}

// HandleCreatePolicyRequest handles the HTTP request to create a network policy from YAML.
func HandleCreatePolicyRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        var policyRequest struct {
            YAML      string `json:"yaml"`
            Namespace string `json:"namespace"`
        }
        if err := json.NewDecoder(r.Body).Decode(&policyRequest); err != nil {
            http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to create Kubernetes client: %v", err), http.StatusInternalServerError)
            return
        }

        networkPolicy, err := YAMLToNetworkPolicy(policyRequest.YAML)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to parse network policy YAML: %v", err), http.StatusBadRequest)
            return
        }

        createdPolicy, err := clientset.NetworkingV1().NetworkPolicies(policyRequest.Namespace).Create(context.Background(), networkPolicy, metav1.CreateOptions{})
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to create network policy: %v", err), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(createdPolicy)
    }
}

func HandleVisualizationRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        namespace := r.URL.Query().Get("namespace")

        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, "Failed to create Kubernetes client: "+err.Error(), http.StatusInternalServerError)
            return
        }

        vizData, err := gatherVisualizationData(clientset, namespace)
        if err != nil {
            http.Error(w, "Failed to gather visualization data: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(vizData); err != nil {
            http.Error(w, "Failed to encode visualization data: "+err.Error(), http.StatusInternalServerError)
        }
    }
}

// HandlePolicyYAMLRequest handles the HTTP request for serving the YAML of a network policy.
func HandlePolicyYAMLRequest(kubeconfigPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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
        clientset, err := GetClientset(kubeconfigPath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        yamlData, err := getNetworkPolicyYAML(clientset, namespace, policyName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/x-yaml")
        w.Write([]byte(yamlData))
    }
}
