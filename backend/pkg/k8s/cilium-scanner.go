package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Use lipgloss for neat tables in CLI
var HeaderAboveTableStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("10")).
	PaddingLeft(0).
	PaddingRight(0).
	MarginBottom(1)

var FoundPolicyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("10")).
	Align(lipgloss.Center).
	PaddingLeft(0).
	PaddingRight(4).
	MarginTop(1).
	MarginBottom(1)

var PoliciesNotApplyingHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("6")).
	Align(lipgloss.Center).
	PaddingLeft(4).
	PaddingRight(4).
	MarginTop(1).
	MarginBottom(1)

var (
	HeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Align(lipgloss.Center)
	EvenRowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	OddRowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
)

func createPodsTable(podsInfo [][]string) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers("Namespace", "Pod Name", "IP Address")

	for _, row := range podsInfo {
		formattedRow := make([]string, 3)
		for i := 0; i < 3; i++ {
			if i < len(row) {
				formattedRow[i] = row[i]
			} else {
				formattedRow[i] = "N/A"
			}
		}

		t.Row(formattedRow[0], formattedRow[1], formattedRow[2])
	}

	return t.String()
}

func createPoliciesTable(policiesInfo [][]string) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers("Policy Name")

	for _, row := range policiesInfo {
		t.Row(row...)
	}

	return t.String()
}

// GetCiliumDynamicClient returns a dynamic interface to query for Cilium policies
func GetCiliumDynamicClient(kubeconfigPath string) (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create k8s client config: %s", err)
		}
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return dynamicClient, nil
}

// initializeCiliumClients creates and returns initialized dynamic and Kubernetes clientsets.
func initializeCiliumClients(kubeconfigPath string) (dynamic.Interface, *kubernetes.Clientset, error) {
	dynamicClient, err := GetCiliumDynamicClient(kubeconfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating dynamic Kubernetes client: %s", err)
	}
	if dynamicClient == nil {
		return nil, nil, fmt.Errorf("failed to create dynamic client: client is nil")
	}

	clientset, err := GetClientset(kubeconfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating Kubernetes clientset: %s", err)
	}
	if clientset == nil {
		return nil, nil, fmt.Errorf("failed to create clientset: clientset is nil")
	}

	return dynamicClient, clientset, nil
}

// fetchCiliumPolicies fetches all Cilium network policies within the specified namespace.
func fetchCiliumPolicies(dynamicClient dynamic.Interface, nsName string, writer *bufio.Writer) ([]*unstructured.Unstructured, bool, error) {
	ciliumNPResource := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v2",
		Resource: "ciliumnetworkpolicies",
	}
	policies, err := dynamicClient.Resource(ciliumNPResource).Namespace(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		printToBoth(writer, fmt.Sprintf("Error listing Cilium network policies in namespace %s: %s\n", nsName, err))
		return nil, false, fmt.Errorf("error listing Cilium network policies: %w", err)
	}

	var unstructuredPolicies []*unstructured.Unstructured
	hasDenyAll := false
	for i := range policies.Items {
		policy := &policies.Items[i]
		unstructuredPolicies = append(unstructuredPolicies, policy)
		if IsDefaultDenyAllCiliumPolicy(*policy) {
			hasDenyAll = true
		}
	}

	return unstructuredPolicies, hasDenyAll, nil
}

// helper function to ensure pods are not added to list multiple times
func addUniquePodDetail(podDetails []string, detail string) []string {
    for _, d := range podDetails {
        if d == detail {
            return podDetails // pod already in the list.
        }
    }
    return append(podDetails, detail) // add pod if its not in list
}

// determinePodCoverage identifies unprotected pods in a namespace based on the fetched Cilium policies.
func determinePodCoverage(clientset *kubernetes.Clientset, nsName string, policies []*unstructured.Unstructured, hasDenyAll bool, writer *bufio.Writer) ([]string, error) {
	unprotectedPods := []string{}

	pods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		printToBoth(writer, fmt.Sprintf("Error listing all pods in namespace %s: %s\n", nsName, err))
		return nil, fmt.Errorf("error listing all pods: %w", err)
	}

	for _, pod := range pods.Items {
		// Skip pods that are not in running state
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		podIdentifier := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
        if _, exists := globallyProtectedPods[podIdentifier]; !exists {
            if !IsPodProtected(writer, clientset, pod, policies, hasDenyAll, globallyProtectedPods) {
                unprotectedPodDetails := fmt.Sprintf("%s %s %s", pod.Namespace, pod.Name, pod.Status.PodIP)
                unprotectedPods = addUniquePodDetail(unprotectedPods, unprotectedPodDetails)
            } else {
                globallyProtectedPods[podIdentifier] = struct{}{} // Mark the pod as protected globally
            }
        }
    }

	return unprotectedPods, nil
}

// processNamespacePoliciesCilium processes Cilium network policies for a given namespace to identify unprotected pods.
func processNamespacePoliciesCilium(dynamicClient dynamic.Interface, clientset *kubernetes.Clientset, nsName string, writer *bufio.Writer, scanResult *ScanResult, dryRun bool, isCLI bool) error {
	ciliumPolicies, hasDenyAll, err := fetchCiliumPolicies(dynamicClient, nsName, writer)
	if err != nil {
		return err
	}

	unprotectedPods, err := determinePodCoverage(clientset, nsName, ciliumPolicies, hasDenyAll, writer)
	if err != nil {
		return err
	}

	if len(unprotectedPods) > 0 {
		// Add unprotected pods to scan results for visibility
		scanResult.UnprotectedPods = append(scanResult.UnprotectedPods, unprotectedPods...)

		if isCLI && !dryRun {
			return handleCLIInteractionsCilium(nsName, unprotectedPods, dynamicClient, writer, scanResult, dryRun)
		} else {
			displayUnprotectedPods(nsName, unprotectedPods, writer)
		}
	}

	return nil
}

func handleCLIInteractionsCilium(nsName string, unprotectedPods []string, dynamicClient dynamic.Interface, writer *bufio.Writer, scanResult *ScanResult, dryRun bool) error {
	unprotectedPodDetails := make([][]string, len(unprotectedPods))
	for i, podDetails := range unprotectedPods {
		unprotectedPodDetails[i] = strings.Fields(podDetails)
	}

	tableOutput := createPodsTable(unprotectedPodDetails)
	headerText := fmt.Sprintf("Unprotected pods found in namespace %s:", nsName)
	styledHeaderText := HeaderStyle.Render(headerText)
	printToBoth(writer, styledHeaderText+"\n"+tableOutput+"\n")

	if !dryRun {
		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Do you want to add a default deny all Cilium network policy to the namespace %s?", nsName),
		}
		if err := survey.AskOne(prompt, &confirm, nil); err != nil {
			return fmt.Errorf("failed to prompt for policy application: %s", err)
		}

		if confirm {
			if err := CreateAndApplyDefaultDenyCiliumPolicy(nsName, dynamicClient); err != nil {
				return fmt.Errorf("failed to apply default deny Cilium policy in namespace %s: %s", nsName, err)
			}
			fmt.Printf("Applied default deny Cilium policy in namespace %s\n", nsName)
			scanResult.PolicyChangesMade = true
		} else {
			scanResult.UserDeniedPolicies = true
		}
	}

	return nil
}

// SelectCiliumNamespaces selects namespaces for scanning based on the input criteria
func SelectCiliumNamespaces(clientset *kubernetes.Clientset, specificNamespace string) ([]string, error) {
	var namespaces []string
	if specificNamespace != "" {
		// Check if the specified namespace exists
		_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), specificNamespace, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return nil, fmt.Errorf("namespace %s does not exist", specificNamespace)
			}
			return nil, fmt.Errorf("error checking namespace %s: %v", specificNamespace, err)
		}
		namespaces = append(namespaces, specificNamespace)
	} else {
		// List all namespaces and filter out system namespaces
		nsList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("error listing namespaces: %v", err)
		}
		for _, ns := range nsList.Items {
			if !IsSystemNamespace(ns.Name) {
				namespaces = append(namespaces, ns.Name)
			}
		}
	}
	return namespaces, nil
}

var hasStartedCiliumScan bool = false
var globallyProtectedPods = make(map[string]struct{})

// ScanCiliumNetworkPolicies scans namespaces for Cilium network policies
func ScanCiliumNetworkPolicies(specificNamespace string, dryRun bool, returnResult bool, isCLI bool, printScore bool, printMessages bool, kubeconfigPath string) (*ScanResult, error) {
	var output bytes.Buffer

	unprotectedPodsCount := 0
	scanResult := new(ScanResult)

	writer := bufio.NewWriter(&output)

	dynamicClient, clientset, err := initializeCiliumClients(kubeconfigPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Check if a specific namespace is provided
	var namespacesToScan []string
	namespacesToScan, err = SelectCiliumNamespaces(clientset, specificNamespace)
	if err != nil {
		return nil, err
	}

	missingPoliciesOrUncoveredPods := false
	userDeniedPolicyApplication := false

	if isCLI && !hasStartedCiliumScan {
		fmt.Println("Policy type: Cilium")
		hasStartedCiliumScan = true
	}

	// Process each namespace for policies and unprotected pods
	for _, nsName := range namespacesToScan {
		scanResult.UnprotectedPods = []string{}
		if err := processNamespacePoliciesCilium(dynamicClient, clientset, nsName, writer, scanResult, dryRun, isCLI); err != nil {
			return nil, err
		}
		unprotectedPodsCount += len(scanResult.UnprotectedPods)

		if len(scanResult.UnprotectedPods) > 0 {
			missingPoliciesOrUncoveredPods = true
		}
	}

	writer.Flush()
	if output.Len() > 0 {
		handleOutputAndPromptsCilium(writer, &output)
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

	hasStartedCiliumScan = false
	return scanResult, nil
}

func handleOutputAndPromptsCilium(writer *bufio.Writer, output *bytes.Buffer) {
	saveToFile := false
	prompt := &survey.Confirm{
		Message: "Do you want to save the output to netfetch-cilium.txt?",
	}
	survey.AskOne(prompt, &saveToFile, nil)

	if saveToFile {
		err := os.WriteFile("netfetch-cilium.txt", output.Bytes(), 0644)
		if err != nil {
			errorFileMsg := fmt.Sprintf("Error writing to file: %s\n", err)
			printToBoth(writer, errorFileMsg)
		} else {
			printToBoth(writer, "Output file created: netfetch-cilium.txt\n")
		}
	} else {
		printToBoth(writer, "Output file not created.\n")
	}
}

// fetchCiliumClusterwidePolicies retrieves all Cilium Clusterwide Network Policies using a dynamic client
func fetchCiliumClusterwidePolicies(dynamicClient dynamic.Interface) ([]*unstructured.Unstructured, error) {
	ciliumCCNPResource := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v2",
		Resource: "ciliumclusterwidenetworkpolicies",
	}

	policies, err := dynamicClient.Resource(ciliumCCNPResource).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing CiliumClusterwideNetworkPolicies: %v", err)
	}

	var unstructuredPolicies []*unstructured.Unstructured
	policyMap := make(map[string]bool)

	for i := range policies.Items {
		policy := policies.Items[i]
		policyName := policy.GetName()
		if _, exists := policyMap[policyName]; !exists {
			policyMap[policyName] = true
			unstructuredPolicies = append(unstructuredPolicies, &policy)
		}
	}

	return unstructuredPolicies, nil
}

// reportDetectedPolicies prints the detected policies to the writer
func reportClusterwideDetectedPolicies(unstructuredPolicies []*unstructured.Unstructured, writer *bufio.Writer, isCLI bool) {
	if isCLI {
		if len(unstructuredPolicies) == 0 {
			printToBoth(writer, "No cluster wide policies found.\n")
		} else {
			for _, policy := range unstructuredPolicies {
				policyName, _, _ := unstructured.NestedString(policy.UnstructuredContent(), "metadata", "name")
				printToBoth(writer, "- "+policyName+"\n")
			}
		}
	}
}

func handleClusterwideCLIInteractions(writer *bufio.Writer, dynamicClient dynamic.Interface, scanResult *ScanResult, appliesToEntireCluster bool, partialDenyAllFound bool, defaultDenyAllFound bool, partialDenyAllPolicies []string, isCLI bool, dryRun bool) error {
	if !appliesToEntireCluster {
		var promptForPolicyCreation bool

		if partialDenyAllFound && defaultDenyAllFound {
			// Display partial policies
			policiesForTable := make([][]string, 0)
			for _, pName := range partialDenyAllPolicies {
				policiesForTable = append(policiesForTable, []string{pName})
			}

			tableOutput := createPoliciesTable(policiesForTable)
			partialPoliciesHeader := HeaderStyle.Render("Cluster wide policies in effect:")
			printToBoth(writer, partialPoliciesHeader+"\n"+tableOutput+"\n")
			promptForPolicyCreation = true
		} else if !defaultDenyAllFound {
			promptForPolicyCreation = true
		}

		if promptForPolicyCreation && isCLI && !dryRun {
			// Prompt to create a default deny-all policy
			createPolicy := false
			prompt := &survey.Confirm{
				Message: "Do you want to create a cluster wide default deny all cilium network policy?",
			}
			if err := survey.AskOne(prompt, &createPolicy, nil); err != nil {
				return fmt.Errorf("failed to prompt for policy application: %s", err)
			}

			if createPolicy {
				if err := CreateAndApplyDefaultDenyCiliumClusterwidePolicy(dynamicClient); err != nil {
					return fmt.Errorf("failed to apply default deny Cilium clusterwide policy: %s", err)
				}
				printToBoth(writer, "\nApplied cluster wide default deny cilium policy\n")
				scanResult.PolicyChangesMade = true
			} else {
				scanResult.UserDeniedPolicies = true
			}
		}
	}

	return nil
}

// checkPodProtection checks each pod against the given policies to determine if it's protected.
func checkPodProtection(clientset *kubernetes.Clientset, unstructuredPolicies []*unstructured.Unstructured, appliesToEntireCluster bool, writer *bufio.Writer) ([]string, error) {
	unprotectedPods := []string{}
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		printToBoth(writer, fmt.Sprintf("Error listing pods: %v\n", err))
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	for _, pod := range pods.Items {
		if !IsSystemNamespace(pod.Namespace) {
			if IsPodProtected(writer, clientset, pod, unstructuredPolicies, appliesToEntireCluster, globallyProtectedPods) {
				podIdentifier := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
				globallyProtectedPods[podIdentifier] = struct{}{}
			} else {
				unprotectedPodDetails := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
				unprotectedPods = append(unprotectedPods, unprotectedPodDetails)
			}
		}
	}
	return unprotectedPods, nil
}

// analyzeClusterwidePolicies processes the list of policies and categorizes them.
func analyzeClusterwidePolicies(unstructuredPolicies []*unstructured.Unstructured) (bool, bool, []string, bool) {
	var defaultDenyAllFound, appliesToEntireCluster, partialDenyAllFound bool
	var partialDenyAllPolicies []string // To hold names of policies that don't apply to the entire cluster

	for _, policy := range unstructuredPolicies {
		isDenyAll, isClusterWide := IsDefaultDenyAllCiliumClusterwidePolicy(*policy)
		if isDenyAll {
			defaultDenyAllFound = true
			if isClusterWide {
				appliesToEntireCluster = true
			} else {
				partialDenyAllFound = true
				policyName, _, _ := unstructured.NestedString(policy.UnstructuredContent(), "metadata", "name")
				partialDenyAllPolicies = append(partialDenyAllPolicies, policyName)
			}
		}
	}
	return defaultDenyAllFound, appliesToEntireCluster, partialDenyAllPolicies, partialDenyAllFound
}

// reportPodProtectionStatus reports the protection status of pods after a scan.
func reportPodProtectionStatus(writer *bufio.Writer, unprotectedPods []string) {
	if len(unprotectedPods) > 0 {
		printToBoth(writer, fmt.Sprintf("Found %d pods not targeted by a cluster wide policy. The namespaced scan will be initiated..\n", len(unprotectedPods)))
	} else {
		printToBoth(writer, "All pods are protected by cluster wide policies.\n")
	}
}

// ScanCiliumClusterwideNetworkPolicies scans the cluster for Cilium Clusterwide Network Policies
func ScanCiliumClusterwideNetworkPolicies(dynamicClient dynamic.Interface, printMessages bool, dryRun bool, isCLI bool, kubeconfigPath string) (*ScanResult, error) {
	// Buffer and writer setup to capture output for both console and file.
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	// Check for a valid dynamic client
	if dynamicClient == nil {
		fmt.Println("Failed to create dynamic client: client is nil")
		return nil, fmt.Errorf("failed to create dynamic client: client is nil")
	}

	dynamicClient, clientset, err := initializeCiliumClients(kubeconfigPath)
	if err != nil {
		fmt.Println("Error initializing clients:", err)
		return nil, err
	}

	unstructuredPolicies, err := fetchCiliumClusterwidePolicies(dynamicClient)
	if err != nil {
		printToBoth(writer, fmt.Sprintf("Error fetching Cilium Clusterwide Network Policies: %s\n", err))
		return nil, err
	}

	if isCLI && !hasStartedCiliumScan {
		fmt.Println("Policy type: Cilium")
		hasStartedCiliumScan = true
	}

	// Report the detected policies
	reportClusterwideDetectedPolicies(unstructuredPolicies, writer, isCLI)

	// Initialize the scan result
	scanResult := &ScanResult{
		NamespacesScanned:  []string{"cluster-wide"},
		DeniedNamespaces:   []string{},
		UnprotectedPods:    []string{},
		PolicyChangesMade:  false,
		UserDeniedPolicies: false,
		AllPodsProtected:   false,
		HasDenyAll:         []string{},
		Score:              50, // or some initial value
	}

	defaultDenyAllFound, appliesToEntireCluster, partialDenyAllPolicies, partialDenyAllFound := analyzeClusterwidePolicies(unstructuredPolicies)

	// Handle CLI interactions for policies
	err = handleClusterwideCLIInteractions(writer, dynamicClient, scanResult, appliesToEntireCluster, partialDenyAllFound, defaultDenyAllFound, partialDenyAllPolicies, isCLI, dryRun)
	if err != nil {
		return nil, err
	}

	// Check pod protection
	unprotectedPods, err := checkPodProtection(clientset, unstructuredPolicies, appliesToEntireCluster, writer)
	if err != nil {
		return nil, err
	}
	scanResult.UnprotectedPods = unprotectedPods

	reportPodProtectionStatus(writer, unprotectedPods)

	if printMessages {
		printToBoth(writer, "\nCluster wide cilium network policy scan completed!\n")
	}

	writer.Flush()
	if output.Len() > 0 {
		handleOutputAndPromptsClusterwideCilium(writer, &output)
	}

	hasStartedCiliumScan = true
	return scanResult, nil
}

func handleOutputAndPromptsClusterwideCilium(writer *bufio.Writer, output *bytes.Buffer) {
	saveToFile := false
	prompt := &survey.Confirm{
		Message: "Do you want to save the output to netfetch-clusterwide-cilium.txt?",
	}
	survey.AskOne(prompt, &saveToFile, nil)

	if saveToFile {
		err := os.WriteFile("netfetch-clusterwide-cilium.txt", output.Bytes(), 0644)
		if err != nil {
			printToBoth(writer, fmt.Sprintf("Error writing to file: %s\n", err))
		} else {
			printToBoth(writer, "Output file created: netfetch-clusterwide-cilium.txt\n")
		}
	} else {
		printToBoth(writer, "Output file not created.\n")
	}
}

func isProtectedByDefaultDeny(policy *unstructured.Unstructured, globallyProtectedPods map[string]struct{}, podIdentifier string) bool {
	_, appliesToEntireCluster := IsDefaultDenyAllCiliumClusterwidePolicy(*policy)
	if appliesToEntireCluster {
		globallyProtectedPods[podIdentifier] = struct{}{}
		return true
	}
	return false
}

func isProtectedByLabelMatch(policies []*unstructured.Unstructured, pod corev1.Pod, globallyProtectedPods map[string]struct{}, podIdentifier string) bool {
	for _, policy := range policies {
		endpointSelector, _, _ := unstructured.NestedMap(policy.UnstructuredContent(), "endpointSelector", "matchLabels")
		if MatchesLabels(pod.Labels, endpointSelector) {
			ingress, foundIngress, _ := unstructured.NestedSlice(policy.UnstructuredContent(), "ingress")
			egress, foundEgress, _ := unstructured.NestedSlice(policy.UnstructuredContent(), "egress")

			// Check for deny-all conditions based on empty ingress/egress
			if (foundIngress && (IsEmptyOrOnlyContainsEmptyObjects(ingress) || IsSpecificallyEmpty(ingress))) ||
				(foundEgress && (IsEmptyOrOnlyContainsEmptyObjects(egress) || IsSpecificallyEmpty(egress))) {
				globallyProtectedPods[podIdentifier] = struct{}{}
				return true
			}

			// Additional check for non-empty specific rules
			if foundIngress && !IsEmptyOrOnlyContainsEmptyObjects(ingress) || foundEgress && !IsEmptyOrOnlyContainsEmptyObjects(egress) {
				globallyProtectedPods[podIdentifier] = struct{}{}
				return true
			}
		}
	}
	return false
}

func IsPodProtected(writer *bufio.Writer, clientset *kubernetes.Clientset, pod corev1.Pod, policies []*unstructured.Unstructured, defaultDenyAllExists bool, globallyProtectedPods map[string]struct{}) bool {
	podIdentifier := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)

	// Immediate return if already protected
	if _, protected := globallyProtectedPods[podIdentifier]; protected {
		return true
	}

	// Apply default deny-all if it exists
	if defaultDenyAllExists {
		globallyProtectedPods[podIdentifier] = struct{}{}
		return true
	}

	// Check each policy for default deny or label match
	for _, policy := range policies {
		if isProtectedByDefaultDeny(policy, globallyProtectedPods, podIdentifier) {
			return true
		}
		if isProtectedByLabelMatch(policies, pod, globallyProtectedPods, podIdentifier) {
			return true
		}
	}

	return false
}

// Check specifically for a slice that only contains a single empty map ({}), representing a default deny.
func IsSpecificallyEmpty(slice []interface{}) bool {
	return len(slice) == 1 && len(slice[0].(map[string]interface{})) == 0
}

// // Placeholder function for future reference
// func entityMatchesPod(entity string, pod corev1.Pod) bool {
// 	switch entity {
// 	case "all":
// 		// All always matches any pod
// 		return true
// 	case "world":
// 		// Determine if the pod communicates with entities outside the cluster
// 		// This might involve checking the pod's networking configuration, labels, or annotations
// 		// Placeholder logic: return false for now
// 		return false
// 	case "host":
// 		// Check if the pod is using host networking
// 		if pod.Spec.HostNetwork {
// 			return true
// 		}
// 		return false
// 	case "remote-node":
// 		// Check if the pod is intended to communicate with a remote node
// 		// This might involve checking node labels, pod's node affinity, or annotations
// 		// Placeholder logic: return false for now
// 		return false
// 	default:
// 		// Unknown entity type, log it, handle as needed
// 		fmt.Printf("Unknown entity type encountered: %s\n", entity)
// 		return false
// 	}
// }

// MatchesLabels checks if the pod's labels match the policy's endpointSelector
func MatchesLabels(podLabels map[string]string, policySelector map[string]interface{}) bool {
	for key, value := range policySelector {
		if val, ok := value.(string); ok {
			if podVal, podOk := podLabels[key]; !podOk || podVal != val {
				return false
			}
		} else {
			fmt.Printf("Policy label value %v is not a string\n", value)
			return false
		}
	}
	return true
}

// ConvertEndpointToSelector converts the endpointSelector from a CiliumNetworkPolicy to a label selector string.
func ConvertEndpointToSelector(endpointSelector map[string]interface{}) (string, error) {
	matchLabels, found := endpointSelector["matchLabels"].(map[string]interface{})
	if !found || len(matchLabels) == 0 {
		return "", nil
	}

	var selectorParts []string
	for key, val := range matchLabels {
		if strVal, ok := val.(string); ok {
			selectorParts = append(selectorParts, fmt.Sprintf("%s=%s", key, strVal))
		} else {
			return "", fmt.Errorf("value for %s in matchLabels is not a string", key)
		}
	}

	labelSelector := strings.Join(selectorParts, ",")
	return labelSelector, nil
}

// CreateAndApplyDefaultDenyCiliumClusterwidePolicy creates and applies a default deny all network policy for Cilium at the cluster level.
func CreateAndApplyDefaultDenyCiliumClusterwidePolicy(dynamicClient dynamic.Interface) error {
	// Construct the policy
	denyAllPolicy := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cilium.io/v2",
			"kind":       "CiliumClusterwideNetworkPolicy",
			"metadata": map[string]interface{}{
				"name": "clusterwide-default-deny-all",
			},
			"spec": map[string]interface{}{
				"endpointSelector": map[string]interface{}{
					"matchLabels": map[string]string{},
				},
				"ingress": []interface{}{},
				"egress":  []interface{}{},
			},
		},
	}

	// Set the GVR
	ciliumCCNPResource := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v2",
		Resource: "ciliumclusterwidenetworkpolicies",
	}

	// Create the policy
	_, err := dynamicClient.Resource(ciliumCCNPResource).Create(context.Background(), denyAllPolicy, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create default deny all CiliumClusterwideNetworkPolicy: %v", err)
	}

	return nil
}

// IsDefaultDenyAllCiliumClusterwidePolicy checks if a single CiliumClusterwideNetworkPolicy is a default deny-all policy
func IsDefaultDenyAllCiliumClusterwidePolicy(policyUnstructured unstructured.Unstructured) (bool, bool) {
	spec, found := policyUnstructured.UnstructuredContent()["spec"].(map[string]interface{})
	if !found {
		fmt.Println("Spec not found in policy.")
		return false, false
	}

	ingress, ingressFound := spec["ingress"].([]interface{})
	egress, egressFound := spec["egress"].([]interface{})

	// Check if ingress and egress are nil or only contain empty objects
	isIngressEmpty := !ingressFound || IsEmptyOrOnlyContainsEmptyObjects(ingress)
	isEgressEmpty := !egressFound || IsEmptyOrOnlyContainsEmptyObjects(egress)

	// Check for default deny-all
	denyAll := isIngressEmpty && isEgressEmpty

	// Check if it applies to the entire cluster
	endpointSelector, selectorFound := spec["endpointSelector"].(map[string]interface{})
	appliesToEntireCluster := !selectorFound || (len(endpointSelector) == 0 || isEndpointSelectorEmpty(endpointSelector))

	return denyAll, appliesToEntireCluster
}

// Helper function to check if the ingress/egress slice is empty or only contains empty objects
func IsEmptyOrOnlyContainsEmptyObjects(slice []interface{}) bool {
	if len(slice) == 0 {
		return true
	}
	for _, item := range slice {
		if itemMap, ok := item.(map[string]interface{}); !ok || len(itemMap) != 0 {
			return false
		}
	}
	return true
}

// Helper function to check if the endpointSelector is effectively empty
func isEndpointSelectorEmpty(selector map[string]interface{}) bool {
	matchLabels, found := selector["matchLabels"].(map[string]interface{})
	return !found || len(matchLabels) == 0
}

// CreateAndApplyDefaultDenyCiliumPolicy creates and applies a default deny all network policy for Cilium in the specified namespace.
func CreateAndApplyDefaultDenyCiliumPolicy(namespace string, dynamicClient dynamic.Interface) error {
	// Construct the policy name dynamically
	policyName := namespace + "-cilium-default-deny-all"

	// Define the policy
	denyAllPolicy := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cilium.io/v2",
			"kind":       "CiliumNetworkPolicy",
			"metadata": map[string]interface{}{
				"name":      policyName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"endpointSelector": map[string]interface{}{
					"matchLabels": map[string]string{},
				},
				"ingress": []interface{}{},
				"egress":  []interface{}{},
			},
		},
	}

	// Set the GVR
	ciliumNPResource := schema.GroupVersionResource{
		Group:    "cilium.io",
		Version:  "v2",
		Resource: "ciliumnetworkpolicies",
	}

	// Create the policy
	_, err := dynamicClient.Resource(ciliumNPResource).Namespace(namespace).Create(context.TODO(), denyAllPolicy, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create default deny all CiliumNetworkPolicy: %v", err)
	}

	fmt.Printf("Applied default deny all CiliumNetworkPolicy to namespace %s\n", namespace)
	return nil
}

// HasDefaultDenyAllCiliumPolicy checks if the list of CiliumNetworkPolicies includes a default deny all policy
func HasDefaultDenyAllCiliumPolicy(policies []*unstructured.Unstructured) bool {
	for _, policy := range policies {
		if IsDefaultDenyAllCiliumPolicy(*policy) {
			return true
		}
	}
	return false
}

// IsDefaultDenyAllCiliumPolicy checks if a single Cilium policy is a default deny-all policy
func IsDefaultDenyAllCiliumPolicy(policyUnstructured unstructured.Unstructured) bool {
	spec, found := policyUnstructured.UnstructuredContent()["spec"].(map[string]interface{})
	if !found {
		return false
	}

	ingress, ingressFound := spec["ingress"]
	egress, egressFound := spec["egress"]
	return (!ingressFound || len(ingress.([]interface{})) == 0) && (!egressFound || len(egress.([]interface{})) == 0)
}
