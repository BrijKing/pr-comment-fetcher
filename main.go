package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Webhook called -----")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	print(body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v\n", err)
		http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
		return
	}
	log.Printf("Raw Request Body start heare ====================================:\n%s\n", prettyJSON.Bytes())

	// Parse the JSON payload
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	for i, v := range payload {
		println(i, v)
	}

	if comment, ok := payload["comment"]; ok {
		if commentMap, ok := comment.(map[string]interface{}); ok {
			if commentBody, ok := commentMap["body"].(string); ok {
				log.Printf("Comment Body: %s\n", commentBody)
			} else {
				log.Println("Comment body not found or not a string")
			}
		} else {
			log.Println("Comment data is not in expected format")
		}
	} else {
		log.Println("No comment data found in payload")
	}

	w.WriteHeader(http.StatusOK)
}
