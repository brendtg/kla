package cluster

import (
	"encoding/json"
	"kla/internal/model"
	"math"
	"strings"
)

// Cluster represents a group of similar log entries
type Cluster struct {
	Group          int    `json:"group"`
	Representative string `json:"representative"`
	Size           int    `json:"size"`
}

// ComputeCosineSimilarity computes the cosine similarity between two feature vectors
func ComputeCosineSimilarity(vec1, vec2 map[string]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for key, value1 := range vec1 {
		value2 := vec2[key]
		dotProduct += value1 * value2
		magnitude1 += value1 * value1
		magnitude2 += value2 * value2
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

// ConvertToFeatureVector converts a LogEntry to a feature vector
func ConvertToFeatureVector(log model.LogEntry) map[string]float64 {
	features := make(map[string]float64)
	for _, word := range strings.Fields(log.Message) {
		features[word]++
	}
	return features
}

// ClusterLogs clusters the logs based on cosine similarity
func ClusterLogs(logs []model.LogEntry, threshold float64) [][]model.LogEntry {
	var clusters [][]model.LogEntry
	used := make([]bool, len(logs))

	for i := 0; i < len(logs); i++ {
		if used[i] {
			continue
		}
		cluster := []model.LogEntry{logs[i]}
		used[i] = true

		for j := i + 1; j < len(logs); j++ {
			if used[j] {
				continue
			}
			similarity := ComputeCosineSimilarity(ConvertToFeatureVector(logs[i]), ConvertToFeatureVector(logs[j]))
			if similarity >= threshold {
				cluster = append(cluster, logs[j])
				used[j] = true
			}
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}

// PrepareJSONOutput prepares the JSON output for the clusters
func PrepareJSONOutput(clusters [][]model.LogEntry) (string, error) {
	var jsonClusters []Cluster
	for i, cluster := range clusters {
		if len(cluster) > 0 {
			jsonClusters = append(jsonClusters, Cluster{
				Group:          i + 1,
				Representative: cluster[0].Message,
				Size:           len(cluster),
			})
		}
	}

	jsonOutput, err := json.MarshalIndent(jsonClusters, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonOutput), nil
}
