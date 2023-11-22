package k8s

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Scan namespace for pods function
func ScanNamespace(namespace string) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %s\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %s\n", err)
		return
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %s\n", err)
		return
	}

	for _, pod := range pods.Items {
		fmt.Printf("Name: %s, Namespace: %s\n", pod.Name, pod.Namespace)
	}
}

// ListPods lists all pods in the specified namespace and returns their names
func ListPods(namespace string) ([]string, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Error building kubeconfig: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating Kubernetes client: %s", err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error listing pods: %s", err)
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

// SniffPods sniffs network traffic from a list of pods in a namespace
func SniffPods(podNames []string, namespace string) error {
	for _, podName := range podNames {
		fmt.Printf("Sniffing pod: %s\n", podName)
		err := SniffPodTraffic(podName, namespace)
		if err != nil {
			fmt.Printf("Error sniffing pod %s: %s\n", podName, err)
			// Decide whether to continue, return, or handle the error differently
		}
		// Add a delay between sniffs to manage resource usage, if necessary
		time.Sleep(5 * time.Second)
	}
	return nil
}

// SniffPodTraffic runs ksniff on a specific pod
func SniffPodTraffic(podName, namespace string) error {
	cmd := exec.Command("kubectl", "sniff", podName, "-n", namespace, "-o", podName+".pcap")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting network sniffing on pod %s in namespace %s\n", podName, namespace)
	return cmd.Run()
}
