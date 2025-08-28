package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Node struct {
	XMLName  xml.Name `json:"tag"`
	Content  string   `xml:",chardata" json:"content,omitempty"`
	Children []Node   `xml:",any" json:"children,omitempty"`
}

// parse SSML into Node Tree
func ParseSSML(data string) (*Node, error) {
	var root Node
	err := xml.Unmarshal([]byte(data), &root)
	if err != nil {
		return nil, err
	}
	return &root, nil
}

// HTTP handler
func ssmlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	node, err := ParseSSML(string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse SSML: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

func main() {
	http.HandleFunc("/parse", ssmlHandler)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
