package k8s

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Helper function to write to both buffer and standard output
func printToBoth(writer *bufio.Writer, s string) {
	// Print to standard output with ANSI codes
	fmt.Print(s)

	// Write to buffer without ANSI codes
	cleanString := StripANSICodes(s)
	fmt.Fprint(writer, cleanString)
}

// StripANSICodes removes ANSI escape codes from a string
func StripANSICodes(str string) string {
	ansi := regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)
	return ansi.ReplaceAllString(str, "")
}

// Struct to represent scan results in dashboard
type ScanResult struct {
	NamespacesScanned  []string
	DeniedNamespaces   []string
	UnprotectedPods    []string
	PolicyChangesMade  bool
	UserDeniedPolicies bool
	HasDenyAll         []string
	Score              int
	AllPodsProtected   bool
}

// Check if error scanning is related to network issues
func isNetworkError(err error) bool {
	var urlError *url.Error
	var netOpError *net.OpError
	var dnsError *net.DNSError

	if errors.As(err, &urlError) {
		if errors.As(urlError.Err, &netOpError) {
			if errors.As(netOpError.Err, &dnsError) {
				if dnsError.IsNotFound {
					return true
				}
			}
		}
	}
	return false
}

// Initialize client
func InitializeClient() (*kubernetes.Clientset, error) {
	clientset, err := GetClientset()
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %s\n", err)
		return nil, err
	}
	return clientset, nil
}

// Select which namespace to scan
func SelectNamespaces(clientset *kubernetes.Clientset, specificNamespace string) ([]string, error) {
	var namespaces []string
	if specificNamespace != "" {
		_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), specificNamespace, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return nil, fmt.Errorf("namespace %s does not exist", specificNamespace)
			}
			return nil, fmt.Errorf("error checking namespace %s: %w", specificNamespace, err)
		}
		namespaces = append(namespaces, specificNamespace)
	} else {
		nsList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			if isNetworkError(err) {
				return nil, fmt.Errorf("network error while listing namespaces, please check your connection to the Kubernetes cluster: %w", err)
			}
			return nil, fmt.Errorf("error listing namespaces: %w", err)
		}
		for _, ns := range nsList.Items {
			if !IsSystemNamespace(ns.Name) {
				namespaces = append(namespaces, ns.Name)
			}
		}
	}
	return namespaces, nil
}

// promptForPolicyApplication asks the user whether to apply a default deny policy
func promptForPolicyApplication(namespace string, writer *bufio.Writer) bool {
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Do you want to add a default deny all network policy to the namespace %s?", namespace),
	}
	err := survey.AskOne(prompt, &confirm)
	if err != nil {
		fmt.Fprintf(writer, "Error prompting for policy application: %s\n", err)
		return false
	}
	return confirm
}

// Fetches all network policies for a namespace and returns a map of covered pods
func fetchCoveredPods(clientset *kubernetes.Clientset, nsName string, writer *bufio.Writer) (map[string]bool, error) {
	coveredPods := make(map[string]bool)
	policies, err := clientset.NetworkingV1().NetworkPolicies(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		printToBoth(writer, fmt.Sprintf("\nError listing network policies in namespace %s: %s\n", nsName, err))
		return nil, fmt.Errorf("error listing network policies: %w", err)
	}

	for _, policy := range policies.Items {
		selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.PodSelector)
		if err != nil {
			printToBoth(writer, fmt.Sprintf("Error parsing selector for policy %s: %s\n", policy.Name, err))
			continue
		}
		pods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			printToBoth(writer, fmt.Sprintf("Error listing pods for policy %s: %s\n", policy.Name, err))
			continue
		}
		for _, pod := range pods.Items {
			coveredPods[pod.Name] = true
		}
	}
	return coveredPods, nil
}

// Fetches all pods in a namespace and determines which are unprotected
func determineUnprotectedPods(clientset *kubernetes.Clientset, nsName string, coveredPods map[string]bool, writer *bufio.Writer, scanResult *ScanResult) ([]string, error) {
	unprotectedPods := []string{}
	allPods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		printToBoth(writer, fmt.Sprintf("Error listing all pods in namespace %s: %s\n", nsName, err))
		return nil, fmt.Errorf("error listing all pods: %w", err)
	}

	for _, pod := range allPods.Items {
		// Skip pods that are not in running state
		if pod.Status.Phase != v1.PodRunning {
			continue
		}
		if !coveredPods[pod.Name] {
			podDetail := fmt.Sprintf("%s %s %s", nsName, pod.Name, pod.Status.PodIP)
			if !containsPodDetail(scanResult.UnprotectedPods, podDetail) {
				unprotectedPods = append(unprotectedPods, podDetail)
			}
		}
	}
	return unprotectedPods, nil
}

// This function just displays unprotected pods without prompting for any actions
func displayUnprotectedPods(nsName string, unprotectedPods []string, writer *bufio.Writer) {
	if len(unprotectedPods) > 0 {
		headerText := fmt.Sprintf("Unprotected pods found in namespace %s:", nsName)
		styledHeaderText := HeaderStyle.Render(headerText)
		printToBoth(writer, styledHeaderText+"\n")

		podsInfo := make([][]string, len(unprotectedPods))
		for i, podDetail := range unprotectedPods {
			podsInfo[i] = strings.Fields(podDetail)
		}

		tableOutput := createPodsTable(podsInfo)
		printToBoth(writer, tableOutput+"\n")
	}
}

func handleCLIInteractions(nsName string, unprotectedPods []string, writer *bufio.Writer, scanResult *ScanResult) {
	if len(unprotectedPods) > 0 {
		// Header
		headerText := fmt.Sprintf("Unprotected pods found in namespace %s:", nsName)
		styledHeaderText := HeaderStyle.Render(headerText)
		printToBoth(writer, styledHeaderText+"\n")

		// Prepare data for table
		podsInfo := make([][]string, len(unprotectedPods))
		for i, podDetail := range unprotectedPods {
			podsInfo[i] = strings.Fields(podDetail)
		}

		// Create and print table
		tableOutput := createPodsTable(podsInfo)
		printToBoth(writer, tableOutput+"\n")

		// Prompt for applying policies
		if promptForPolicyApplication(nsName, writer) {
			err := createAndApplyDefaultDenyPolicy(nsName)
			if err != nil {
				fmt.Fprintf(writer, "Failed to apply default deny policy in namespace %s: %s\n", nsName, err)
			} else {
				fmt.Fprintf(writer, "Applied default deny policy in namespace %s\n", nsName)
				scanResult.PolicyChangesMade = true
			}
		}
	}
}

func processNamespacePolicies(clientset *kubernetes.Clientset, nsName string, writer *bufio.Writer, isCLI bool, dryRun bool, scanResult *ScanResult) error {
	// Fetch covered pods
	coveredPods, err := fetchCoveredPods(clientset, nsName, writer)
	if err != nil {
		return fmt.Errorf("fetching covered pods failed for namespace %s: %w", nsName, err)
	}

	// Determine unprotected pods
	unprotectedPods, err := determineUnprotectedPods(clientset, nsName, coveredPods, writer, scanResult)
	if err != nil {
		return fmt.Errorf("determining unprotected pods failed for namespace %s: %w", nsName, err)
	}

	// Always add pods to result for visibility
	scanResult.UnprotectedPods = append(scanResult.UnprotectedPods, unprotectedPods...)
	scanResult.DeniedNamespaces = append(scanResult.DeniedNamespaces, nsName)

	// Only handle CLI interactions if it's CLI mode and not a dry run
	if isCLI && !dryRun {
		handleCLIInteractions(nsName, unprotectedPods, writer, scanResult)
	} else if dryRun {
		// If it's a dry run, we just display the data without prompting for any actions
		displayUnprotectedPods(nsName, unprotectedPods, writer)
	}

	return nil
}

var hasStartedNativeScan bool = false

// ScanNetworkPolicies scans namespaces for network policies
func ScanNetworkPolicies(specificNamespace string, dryRun bool, returnResult bool, isCLI bool, printScore bool, printMessages bool) (*ScanResult, error) {
	var output bytes.Buffer
	var namespacesToScan []string

	unprotectedPodsCount := 0
	scanResult := new(ScanResult)

	writer := bufio.NewWriter(&output)

	clientset, err := InitializeClient()
	if err != nil {
		return nil, err
	}

	namespacesToScan, err = SelectNamespaces(clientset, specificNamespace)
	if err != nil {
		return nil, err
	}

	missingPoliciesOrUncoveredPods := false
	userDeniedPolicyApplication := false
	deniedNamespaces := []string{}

	if isCLI && !hasStartedNativeScan {
		fmt.Println("Policy type: Kubernetes")
		hasStartedNativeScan = true
	}

	for _, nsName := range namespacesToScan {
		err := processNamespacePolicies(clientset, nsName, writer, isCLI, dryRun, scanResult)
		if err != nil {
			fmt.Printf("Error processing namespace %s: %v\n", nsName, err)
			continue
		}
		unprotectedPodsCount += len(scanResult.UnprotectedPods)

		// Check if namespace is already marked as denied
		if !contains(deniedNamespaces, nsName) {
			deniedNamespaces = append(deniedNamespaces, nsName)
		}
	}

	writer.Flush()
	if output.Len() > 0 {
		handleOutputAndPrompts(writer, &output)
	}

	score := CalculateScore(!missingPoliciesOrUncoveredPods, !userDeniedPolicyApplication, unprotectedPodsCount)
	scanResult.Score = score

	if printMessages {
		printToBoth(writer, "\nNetfetch scan completed!\n")
	}

	if printScore {
		// Print the final score
		fmt.Printf("\nYour Netfetch security score is: %d/100\n", score)
	}

	hasStartedNativeScan = false
	return scanResult, nil
}

// handleOutputAndPrompts manages saving scan results to a file and outputting
func handleOutputAndPrompts(writer *bufio.Writer, output *bytes.Buffer) {
	saveToFile := false
	prompt := &survey.Confirm{
		Message: "Do you want to save the output to netfetch.txt?",
	}
	survey.AskOne(prompt, &saveToFile, nil)

	if saveToFile {
		err := os.WriteFile("netfetch.txt", output.Bytes(), 0644)
		if err != nil {
			errorFileMsg := fmt.Sprintf("Error writing to file: %s\n", err)
			printToBoth(writer, errorFileMsg)
		} else {
			printToBoth(writer, "Output file created: netfetch.txt\n")
		}
	} else {
		printToBoth(writer, "Output file not created.\n")
	}
}

// Function to create the implicit default deny if missing
func createAndApplyDefaultDenyPolicy(namespace string) error {
	// Initialize Kubernetes client
	clientset, err := GetClientset()
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// Define the network policy
	policyName := namespace + "-default-deny-all"
	policy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyName,
			Namespace: namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
		},
	}

	// Create the policy
	_, err = clientset.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), policy, metav1.CreateOptions{})
	return err
}

// hasDefaultDenyAllPolicy checks if the list of policies includes a default deny all policy
func hasDefaultDenyAllPolicy(policies []networkingv1.NetworkPolicy) bool {
	for _, policy := range policies {
		if isDefaultDenyAllPolicy(policy) {
			return true
		}
	}
	return false
}

// isDefaultDenyAllPolicy checks if a single network policy is a default deny all policy
func isDefaultDenyAllPolicy(policy networkingv1.NetworkPolicy) bool {
	return len(policy.Spec.Ingress) == 0 && len(policy.Spec.Egress) == 0
}

// isSystemNamespace checks if the given namespace is a system namespace
func IsSystemNamespace(namespace string) bool {
	switch namespace {
	case "kube-system", "tigera-operator", "kube-public", "kube-node-lease", "gatekeeper-system", "calico-system":
		return true
	default:
		return false
	}
}

// Scoring logic
func CalculateScore(hasPolicies bool, hasDenyAll bool, unprotectedPodsCount int) int {
    score := 50 // Start with a base score of 50

    if hasDenyAll {
        score += 20 // Add 20 points for having deny-all policies
    } else if !hasPolicies {
        score -= 20 // Subtract 20 points if there are no policies at all
    }

    // Deduct score based on the number of unprotected pods
    score -= unprotectedPodsCount

    if score > 100 {
        score = 100
    } else if score < 1 {
        score = 1
    }

    return score
}

// INTERACTIVE DASHBOARD LOGIC

// handleScanRequest handles the HTTP request for scanning network policies
func HandleScanRequest(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from request, e.g., namespace
	namespace := r.URL.Query().Get("namespace")

	// Perform the scan
	result, err := ScanNetworkPolicies(namespace, false, true, false, false, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleNamespaceListRequest lists all non-system Kubernetes namespaces
func HandleNamespaceListRequest(w http.ResponseWriter, r *http.Request) {
	clientset, err := GetClientset()
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

var (
	isClientInitialized = false
	clientset           *kubernetes.Clientset
)

// GetClientset creates a new Kubernetes clientset
func GetClientset() (*kubernetes.Clientset, error) {
	if isClientInitialized {
		return clientset, nil
	}

	var config *rest.Config
	var err error

	// First try to use the in-cluster configuration
	config, err = rest.InClusterConfig()
	if err != nil {
		fmt.Println("Mode: CLI")

		// Fallback to kubeconfig
		var kubeconfig string
		if kc := os.Getenv("KUBECONFIG"); kc != "" {
			kubeconfig = kc
			fmt.Println("Using KUBECONFIG from environment:", kubeconfig)
		} else {
			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
			fmt.Println("Using default kubeconfig path:", kubeconfig)
			fmt.Printf("\n")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubeconfig path %s: %v", kubeconfig, err)
		}
	} else {
		fmt.Println("Using in-cluster Kubernetes configuration")
	}

	// Create and store the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	isClientInitialized = true
	return clientset, nil
}

func HandleAddPolicyRequest(w http.ResponseWriter, r *http.Request) {
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
	err := createAndApplyDefaultDenyPolicy(req.Namespace)
	if err != nil {
		http.Error(w, "Failed to apply default deny policy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Implicit default deny all network policy successfully added to namespace " + req.Namespace})

	scanResult, err := ScanNetworkPolicies(req.Namespace, false, true, false, false, false)
	if err != nil {
		http.Error(w, "Error re-scanning after applying policy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with updated scan results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scanResult)
}

// contains checks if a string is present in a slice
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// containsPodDetail checks if a pod detail string is present in a slice
func containsPodDetail(slice []string, detail string) bool {
	for _, v := range slice {
		if v == detail {
			return true
		}
	}
	return false
}

// PodInfo holds the desired information from a Pods YAML.
type PodInfo struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Ports     []v1.ContainerPort
}

// Hold the desired info from a Pods ports
type ContainerPortInfo struct {
	Name          string
	ContainerPort int32
	Protocol      v1.Protocol
}

func GetPodInfo(clientset kubernetes.Interface, namespace string) ([]PodInfo, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podInfos []PodInfo
	for _, pod := range pods.Items {
		var containerPorts []v1.ContainerPort
		for _, container := range pod.Spec.Containers {
			containerPorts = append(containerPorts, container.Ports...)
		}

		podInfo := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Labels:    pod.Labels,
			Ports:     containerPorts,
		}
		podInfos = append(podInfos, podInfo)
	}

	return podInfos, nil
}

// YAMLToNetworkPolicy converts a YAML string to a NetworkPolicy object.
func YAMLToNetworkPolicy(yamlContent string) (*networkingv1.NetworkPolicy, error) {
	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer()
	obj, _, err := decoder.Decode([]byte(yamlContent), nil, nil)
	if err != nil {
		return nil, err
	}

	networkPolicy, ok := obj.(*networkingv1.NetworkPolicy)
	if !ok {
		return nil, fmt.Errorf("decoded object is not a NetworkPolicy")
	}

	return networkPolicy, nil
}
