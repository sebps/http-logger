package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Response is the structure of the JSON response (used in mirror mode)
type Response struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Body    string `json:"body"`
	Headers string `json:"headers"`
	Params  string `json:"params"`
}

var mirrorMode bool

func main() {
	// Custom help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `HTTP Logger

This tool starts an HTTP server that logs all incoming requests.
By default, it returns "ok" with status 200.
Enable mirror mode to return the full request as JSON.

Usage:
  http-logger [options]

Options:
  -p, --port <port>    Port to run the server on (default: 8080)
  -m, --mirror         Enable mirror mode (responds with full request details)
  -h, --help           Show this help message and exit
`)
	}

	// Flags
	port := flag.String("port", "8080", "Port to run the HTTP server on")
	flag.StringVar(port, "p", "8080", "Port to run the HTTP server on (shorthand)")
	flag.BoolVar(&mirrorMode, "mirror", false, "Enable mirror mode")
	flag.BoolVar(&mirrorMode, "m", false, "Enable mirror mode (shorthand)")
	help := flag.Bool("h", false, "Show help")
	flag.Parse()

	// Manual --help check
	for _, arg := range os.Args[1:] {
		if arg == "--help" {
			flag.Usage()
			os.Exit(0)
		}
	}
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	address := ":" + *port
	http.HandleFunc("/", rootHandler)

	log.Printf("🚀 HTTP Logger starting on http://localhost%s (mirror mode: %v)", address, mirrorMode)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("📥 %s request on path: %s", r.Method, r.URL.Path)

	// Mirror mode logic
	w.Header().Set("Content-Type", "application/json")

	// Read body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	body := string(bodyBytes)

	// Read headers
	headers := []string{}
	for k, v := range r.Header {
		headers = append(headers, fmt.Sprintf("%s: %s", k, strings.Join(v, ",")))
	}
	headersStr := strings.Join(headers, "\n")

	// Read query parameters
	params := []string{}
	for k, v := range r.URL.Query() {
		params = append(params, fmt.Sprintf("%s=%s", k, strings.Join(v, ",")))
	}
	paramsStr := strings.Join(params, "&")

	// Log full request
	log.Printf("🔎 Headers:\n%s", headersStr)
	log.Printf("🔎 Query Parameters: %s", paramsStr)
	log.Printf("🔎 Body: %s", body)

	if !mirrorMode {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	resp := Response{
		Method:  r.Method,
		Path:    r.URL.Path,
		Body:    body,
		Headers: headersStr,
		Params:  paramsStr,
	}

	if r.Method != http.MethodHead {
		json.NewEncoder(w).Encode(resp)
	}
}
