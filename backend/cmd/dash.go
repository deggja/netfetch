package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/charmbracelet/lipgloss"
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
		port, _ := cmd.Flags().GetString("port")
		startDashboardServer(port, kubeconfigPath)
	},
}

func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func startDashboardServer(port string, kubeconfigPath string) {
	// Verify connection to cluster or throw error
	clientset, err := k8s.GetClientset(kubeconfigPath)
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
			allowedOrigins := []string{"http://localhost:" + port, "https://" + host}
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
	http.HandleFunc("/scan", k8s.HandleScanRequest(kubeconfigPath))
	http.HandleFunc("/namespaces", k8s.HandleNamespaceListRequest(kubeconfigPath))
	http.HandleFunc("/add-policy", k8s.HandleAddPolicyRequest(kubeconfigPath))
	http.HandleFunc("/create-policy", k8s.HandleCreatePolicyRequest(kubeconfigPath))
	http.HandleFunc("/namespaces-with-policies", k8s.HandleNamespacesWithPoliciesRequest(kubeconfigPath))
	http.HandleFunc("/namespace-policies", k8s.HandleNamespacePoliciesRequest(kubeconfigPath))
	http.HandleFunc("/visualization", k8s.HandleVisualizationRequest(kubeconfigPath))
	http.HandleFunc("/visualization/cluster", k8s.HandleClusterVisualizationRequest(kubeconfigPath))
	http.HandleFunc("/policy-yaml", k8s.HandlePolicyYAMLRequest(kubeconfigPath))
	http.HandleFunc("/pod-info", k8s.HandlePodInfoRequest(kubeconfigPath))

	// Wrap the default serve mux with the CORS middleware
	handler := c.Handler(http.DefaultServeMux)

	// Start the server
	serverURL := fmt.Sprintf("http://localhost:%s", port)
	startupMessage := HeaderStyle.Render(fmt.Sprintf("Starting dashboard server on %s", serverURL))
	fmt.Println(startupMessage)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

// func dashboardHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check if we are in development mode
// 	isDevelopment := true
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

var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("6")).
	Align(lipgloss.Center).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(4).
	PaddingRight(4).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("99"))

func init() {
	dashCmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file (optional)")
	dashCmd.Flags().StringP("port", "p", "8080", "Port for the interactive dashboard")
	rootCmd.AddCommand(dashCmd)
}