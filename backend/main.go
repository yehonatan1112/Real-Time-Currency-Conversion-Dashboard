package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var ratesChannel = make(chan map[string]float64) // Channel for rate updates

// Struct for holding API response
type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for WebSocket connections
	},
}

// Fetch rates from external API
func fetchRates(apiURL string, channel chan<- map[string]float64) {
	for {
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Println("Error fetching rates:", err)
			time.Sleep(1 * time.Minute) // Retry after 1 minute
			continue
		}
		// Close the response body at the end of the current iteration
		defer resp.Body.Close()

		var rates ExchangeRates
		err = json.NewDecoder(resp.Body).Decode(&rates)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			time.Sleep(1 * time.Minute) // Retry after 1 minute
			continue
		}

		// Send rates to the channel
		channel <- rates.Rates

		time.Sleep(10 * time.Second) // Fetch every 10 seconds
	}
}

// WebSocket handler to send real-time rates
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Get the latest rates from the channel and send them to the client
		rates := <-ratesChannel
		err := conn.WriteJSON(rates) // Send rates as JSON via WebSocket
		if err != nil {
			log.Println("Error sending rates via WebSocket:", err)
			return
		}
	}
}

func main() {
	apiURL := "http://api.exchangeratesapi.io/v1/latest?access_key=95080ca47e214b57c3b268fe2aaa9f12" // Replace with your valid API URL and key

	// Start fetching rates in a separate goroutine
	go fetchRates(apiURL, ratesChannel)

	// Setup WebSocket endpoint
	http.HandleFunc("/live-updates", handleWebSocket)

	// Start the server
	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
