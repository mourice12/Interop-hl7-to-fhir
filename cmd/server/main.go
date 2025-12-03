package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mourice12/hl7-to-fhir/internal/converter"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

func main() {
	http.HandleFunc("/convert", handleConvert)
	http.HandleFunc("/health", handleHealth)

	port := ":8000"
	fmt.Printf("Server starting on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// handleConvert processes hl7 to FHIR conversion
func handleConvert(w http.ResponseWriter, r *http.Request) {
	//Only accept Post
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Read the HL7 Message

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//parse HL7
	msg, err := hl7.Parse(string(body))
	if err != nil {
		http.Error(w, "Error parsing HL7: "+err.Error(), http.StatusBadRequest)
		return
	}

	//Convert to FHIR bundle
	bundle, err := converter.ConvertToBundle(msg)
	if err != nil {
		http.Error(w, "Error Converting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//return JSON response
	w.Header().Set("Content-Type", "application/fhir+json")
	json.NewEncoder(w).Encode(bundle)

}

// handleHealth returns server status
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "healthy"}`))
}
