package cmd

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/deggja/netfetch/backend/statik"

	"github.com/deggja/netfetch/backend/pkg/k8s"
	"github.com/rakyll/statik/fs"
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

func startDashboardServer() {
	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-CSRF-Token"},
	})

	// Set up handlers
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/scan", k8s.HandleScanRequest)
	http.HandleFunc("/namespaces", k8s.HandleNamespaceListRequest)
	http.HandleFunc("/add-policy", k8s.HandleAddPolicyRequest)
	http.HandleFunc("/visualization", k8s.HandleVisualizationRequest)

	// Wrap the default serve mux with the CORS middleware
	handler := c.Handler(http.DefaultServeMux)

	// Start the server
	port := "8080"
	fmt.Printf("Starting dashboard server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// Serve the embedded frontend
	http.FileServer(statikFS).ServeHTTP(w, r)
}
