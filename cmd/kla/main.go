package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	"kla/internal/cluster"
	"kla/internal/cmd"
	"kla/internal/logparser"
	"kla/internal/ui"
)

func main() {
	thresholdPtr := flag.Float64("threshold", 0.5, "similarity threshold for clustering")
	jqFilter := flag.String("jq", ".", "jq filter to apply to the JSON output")
	flag.Parse()

	// Get the list of pods
	podsOutput, err := cmd.ExecuteCommand("kubectl", "get", "pods", "--no-headers", "-o", "custom-columns=:metadata.name")
	if err != nil {
		log.Fatalf("Failed to get pods: %v", err)
	}
	podArray := strings.Fields(podsOutput)

	// Display the menu to select a pod
	selectedPod, err := ui.DisplayMenu("Select a pod:", podArray)
	if err != nil {
		log.Fatalf("Failed to select a pod: %v", err)
	}

	// Get the list of containers for the selected pod
	containersOutput, err := cmd.ExecuteCommand("kubectl", "get", "pod", selectedPod, "-o", "jsonpath={.spec.containers[*].name}")
	if err != nil {
		log.Fatalf("Failed to get containers for pod %s: %v", selectedPod, err)
	}
	containersArray := strings.Fields(containersOutput)

	// Display the menu to select a container
	selectedContainer, err := ui.DisplayMenu("Select a container:", containersArray)
	if err != nil {
		log.Fatalf("Failed to select a container: %v", err)
	}

	// Fetch the logs of the selected container
	logsOutput, err := cmd.ExecuteCommand("kubectl", "logs", selectedPod, "-c", selectedContainer)
	if err != nil {
		log.Fatalf("Failed to get logs for pod %s container %s: %v", selectedPod, selectedContainer, err)
	}

	// Process the logs to find error messages and cluster them
	logs, err := logparser.ParseLogs(logsOutput)
	if err != nil {
		log.Fatalf("Error parsing logs: %v", err)
	}

	clusters := cluster.ClusterLogs(logs, *thresholdPtr)

	// Prepare JSON output
	jsonOutput, err := cluster.PrepareJSONOutput(clusters)
	if err != nil {
		log.Fatalf("Error generating JSON: %v", err)
	}

	// Pass the JSON output to jq
	jqCmd := exec.Command("jq", *jqFilter)
	jqCmd.Stdin = strings.NewReader(string(jsonOutput))
	jqCmd.Stdout = os.Stdout
	jqCmd.Stderr = os.Stderr

	err = jqCmd.Run()
	if err != nil {
		log.Fatalf("Error running jq: %v", err)
	}
}
