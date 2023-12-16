package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "github.com/deggja/netfetch/backend/statik"
	"github.com/rakyll/statik/fs"

	"github.com/deggja/netfetch/backend/pkg/k8s"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

var dashCmd = &cobra.Command{
	Use:   "dash",
	Short: "Launch the Netfetch interactive dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		startDashboardServer()
	},
}

func init() {
	rootCmd.AddCommand(dashCmd)
}

func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func startDashboardServer() {
	// Verify connection to cluster or throw error
	clientset, err := k8s.GetClientset()
	if err != nil {
		log.Fatalf("You are not connected to a Kubernetes cluster. Please connect to a cluster and re-run the command: %v", err)
		return
	}

	_, err = clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("You are not connected to a Kubernetes cluster. Please connect to a cluster and re-run the command: %v", err)
		return
	}

	c := cors.New(cors.Options{
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			// Implement your dynamic origin check here
			host := r.Host // Extract the host from the request
			allowedOrigins := []string{"http://localhost:8081", "https://" + host}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	// Set up handlers
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/scan", k8s.HandleScanRequest)
	http.HandleFunc("/namespaces", k8s.HandleNamespaceListRequest)
	http.HandleFunc("/add-policy", k8s.HandleAddPolicyRequest)
	http.HandleFunc("/create-policy", HandleCreatePolicyRequest)
	http.HandleFunc("/namespaces-with-policies", handleNamespacesWithPoliciesRequest)
	http.HandleFunc("/namespace-policies", handleNamespacePoliciesRequest)
	http.HandleFunc("/visualization", k8s.HandleVisualizationRequest)
	http.HandleFunc("/visualization/cluster", handleClusterVisualizationRequest)
	http.HandleFunc("/policy-yaml", k8s.HandlePolicyYAMLRequest)
	http.HandleFunc("/pod-info", handlePodInfoRequest)

	// Wrap the default serve mux with the CORS middleware
	handler := c.Handler(http.DefaultServeMux)

	// Start the server
	port := "8080"
	fmt.Printf("Starting dashboard server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

// func dashboardHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check if we are in development mode
// 	isDevelopment := true // You can use an environment variable or a config flag to set this
// 	if isDevelopment {
// 		// Redirect to the Vue dev server
// 		vueDevServer := "http://localhost:8081"
// 		http.Redirect(w, r, vueDevServer+r.RequestURI, http.StatusTemporaryRedirect)
// 	} else {
// 		// Serve the embedded frontend using statik
// 		statikFS, err := fs.New()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		http.FileServer(statikFS).ServeHTTP(w, r)
// 	}
// }

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Set cache control headers
	setNoCacheHeaders(w)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// Serve the embedded frontend
	http.FileServer(statikFS).ServeHTTP(w, r)
}

// handleNamespacesWithPoliciesRequest handles the HTTP request for serving a list of namespaces with network policies.
func handleNamespacesWithPoliciesRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	namespaces, err := k8s.GatherNamespacesWithPolicies()
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

// handleNamespacePoliciesRequest handles the HTTP request for serving a list of network policies in a namespace.
func handleNamespacePoliciesRequest(w http.ResponseWriter, r *http.Request) {
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
	clientset, err := k8s.GetClientset()
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

	// Convert the list of network policies to a more simple structure if needed or encode directly
	// For example, you might want to return only the names and some identifiers of the policies

	setNoCacheHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

// handleClusterVisualizationRequest handles the HTTP request for serving cluster-wide visualization data.
func handleClusterVisualizationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call the function to gather cluster-wide visualization data
	clusterVizData, err := k8s.GatherClusterVisualizationData()
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

// handlePodInfoRequest handles the HTTP request for serving pod information.
func handlePodInfoRequest(w http.ResponseWriter, r *http.Request) {
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
	clientset, err := k8s.GetClientset()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Kubernetes client: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch pod information from the specified namespace
	podInfo, err := k8s.GetPodInfo(clientset, namespace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get pod information: %v", err), http.StatusInternalServerError)
		return
	}

	setNoCacheHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(podInfo)
}

// HandleCreatePolicyRequest handles the HTTP request to create a network policy from YAML.
func HandleCreatePolicyRequest(w http.ResponseWriter, r *http.Request) {
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

	clientset, err := k8s.GetClientset()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Kubernetes client: %v", err), http.StatusInternalServerError)
		return
	}

	networkPolicy, err := k8s.YAMLToNetworkPolicy(policyRequest.YAML)
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
