package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Provider struct {
	Name         string
	ResponseTime time.Duration
	ErrorRate    float64
}

type ProviderManager struct {
	providers      []*Provider
	activeProvider *Provider
	mu             sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (pm *ProviderManager) selectBestProvider() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	bestProvider := pm.providers[0]
	for _, provider := range pm.providers[1:] {
		if provider.ResponseTime < bestProvider.ResponseTime && rand.Float64() > provider.ErrorRate {
			bestProvider = provider
		}
	}
	pm.activeProvider = bestProvider
	log.Printf("Switched to provider: %s", bestProvider.Name)
}

func (pm *ProviderManager) setActiveProvider(provider *Provider) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.activeProvider = provider
}

func (pm *ProviderManager) getActiveProvider() *Provider {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.activeProvider
}

func (pm *ProviderManager) simulateProviderResponse(provider *Provider, ch chan string) {
	time.Sleep(provider.ResponseTime)
	if rand.Float64() < provider.ErrorRate {
		ch <- fmt.Sprintf("Error from provider: %s", provider.Name)
	} else {
		ch <- fmt.Sprintf("Response from provider: %s", provider.Name)
	}
}

func reader(pm *ProviderManager, ws *websocket.Conn) {
	for {
		provider := pm.getActiveProvider()
		ch := make(chan string)

		go pm.simulateProviderResponse(provider, ch)

		select {
		case response := <-ch:
			if err := ws.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("WriteMessage error:", err)
				return
			}

			// Evaluate provider performance and decide if a switch is necessary
			if response[:5] == "Error" {
				pm.selectBestProvider()
			}
		case <-time.After(3 * time.Second): // Timeout for provider response
			log.Println("Provider response timeout, switching provider.")
			pm.selectBestProvider()
		}
	}
}

func streamHandler(pm *ProviderManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer ws.Close()

		reader(pm, ws)
	}
}

func main() {
	// Create stub providers
	providers := []*Provider{
		{Name: "ProviderA", ResponseTime: 1 * time.Second, ErrorRate: 0.1},
		{Name: "ProviderB", ResponseTime: 2 * time.Second, ErrorRate: 0.2},
		{Name: "ProviderC", ResponseTime: 3 * time.Second, ErrorRate: 0.3},
	}

	// Initialize the provider manager with providers
	pm := &ProviderManager{providers: providers, activeProvider: providers[1]}

	// Start net/http server
	http.HandleFunc("/stream", streamHandler(pm))
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
