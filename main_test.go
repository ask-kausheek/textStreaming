package main

import (
	"testing"
	"time"
)

// Helper function to create a new ProviderManager with test providers
func newTestProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: []*Provider{
			{Name: "ProviderA", ResponseTime: 1 * time.Second, ErrorRate: 0.1},
			{Name: "ProviderB", ResponseTime: 2 * time.Second, ErrorRate: 0.2},
			{Name: "ProviderC", ResponseTime: 3 * time.Second, ErrorRate: 0.3},
		},
		activeProvider: nil,
	}
}

// Test that the best provider is selected initially
func TestInitialBestProviderSelection(t *testing.T) {
	pm := newTestProviderManager()
	pm.selectBestProvider()
	if pm.activeProvider.Name != "ProviderA" {
		t.Errorf("Expected ProviderA, got %s", pm.activeProvider.Name)
	}
}

// Test switching to a provider with a lower error rate
func TestProviderSwitchOnErrorRate(t *testing.T) {
	pm := newTestProviderManager()
	pm.setActiveProvider(pm.providers[2]) // Set to ProviderC
	pm.selectBestProvider()
	if pm.activeProvider.Name != "ProviderA" {
		t.Errorf("Expected ProviderA, got %s", pm.activeProvider.Name)
	}
}

// Test switching based on response time
func TestProviderSwitchOnResponseTime(t *testing.T) {
	pm := newTestProviderManager()
	// Simulate slow response from current provider
	pm.setActiveProvider(&Provider{Name: "ProviderD", ResponseTime: 4 * time.Second, ErrorRate: 0.1})
	pm.selectBestProvider()
	if pm.activeProvider.Name != "ProviderA" {
		t.Errorf("Expected ProviderA, got %s", pm.activeProvider.Name)
	}
}

// Test that the active provider doesn't change if it's already the best one
func TestNoSwitchIfBestProvider(t *testing.T) {
	pm := newTestProviderManager()
	pm.setActiveProvider(pm.providers[0]) // Set to ProviderA
	pm.selectBestProvider()
	if pm.activeProvider.Name != "ProviderA" {
		t.Errorf("Expected ProviderA to remain active, got %s", pm.activeProvider.Name)
	}
}
