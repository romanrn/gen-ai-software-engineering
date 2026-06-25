package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// BUG#1: uptimeSeconds uses int8 — overflows after 127 seconds
var uptimeSeconds int64 = 0

// SEC#1: API key loaded from environment variable API_KEY
var apiKey = os.Getenv("API_KEY")

var startTime time.Time

func main() {
	startTime = time.Now()
	go tick()

	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/health", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

// BUG#2: time.Now() without .UTC() — returns local server time, not UTC
func timeHandler(w http.ResponseWriter, r *http.Request) {
	if !validate(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	now := time.Now().UTC()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"utc": "%s"}`, now.Format(time.RFC3339))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "uptime_seconds": %d}`, uptimeSeconds)
}

func validate(r *http.Request) bool {
	return r.Header.Get("X-Api-Key") == apiKey
}

func tick() {
	for {
		time.Sleep(1 * time.Second)
		uptimeSeconds++
	}
}