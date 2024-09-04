package main

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ProviderResult struct {
	Output  string
	Metrics ProviderMetrics
}

type ProviderMetrics struct {
	ResponseTime time.Duration
	Accuracy     float64
}

func sendRequestToProvider(provider string, request string) (string, error) {
	// Simulate provider request
	time.Sleep(100 * time.Millisecond) // Simulate response time
	return "Response from " + provider, nil
}

func fetchFromProviders(providers []string, request string) map[string]ProviderResult {
	results := make(map[string]ProviderResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, provider := range providers {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			start := time.Now()
			output, err := sendRequestToProvider(p, request)
			metrics := ProviderMetrics{ResponseTime: time.Since(start)}

			if err != nil {
				metrics.Accuracy = 0 // Example error handling
			} else {
				metrics.Accuracy = 1.0 // Example accuracy
			}

			mu.Lock()
			results[p] = ProviderResult{Output: output, Metrics: metrics}
			mu.Unlock()
		}(provider)
	}
	wg.Wait()
	return results
}

func evaluateQuality(results map[string]ProviderResult) string {
	var bestProvider string
	var bestScore float64

	for provider, result := range results {
		score := result.Metrics.Accuracy // Simplified example
		if score > bestScore {
			bestScore = score
			bestProvider = provider
		}
	}
	return bestProvider
}

func handleRequest(c *gin.Context) {
	request := c.Query("input")
	providers := []string{"ProviderA", "ProviderB", "ProviderC"}

	results := fetchFromProviders(providers, request)
	bestProvider := evaluateQuality(results)

	c.JSON(200, gin.H{
		"provider": bestProvider,
		"response": results[bestProvider].Output,
	})
}

func main() {
	r := gin.Default()
	r.GET("/process", handleRequest)
	r.Run(":8080")
}
