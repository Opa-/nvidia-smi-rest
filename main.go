package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/clbanning/mxj"
)

func main() {
	port := flag.Int("port", 7777, "HTTP server port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		cmd := exec.Command("nvidia-smi", "-x", "-q")
		output, err := cmd.Output()
		if err != nil {
			JSONError(w, "Error executing nvidia-smi command", http.StatusInternalServerError)
			return
		}

		jsonData, err := xmlToJSON(output)
		if err != nil {
			JSONError(w, "Error converting XML data to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonData)))

		_, err = w.Write(jsonData)
		if err != nil {
			JSONError(w, "Error writing response", http.StatusInternalServerError)
			return
		}
	})

	fmt.Printf("Starting server on port: %d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Printf("Error starting HTTP server: %s\n", err)
	}
}

func xmlToJSON(xmlData []byte) ([]byte, error) {
	m, err := mxj.NewMapXml(xmlData)
	if err != nil {
		return nil, err
	}

	jsonData, err := m.Json()
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

type errorResponse struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}
