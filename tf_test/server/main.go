package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Tag struct {
	Tag string `json:"tag"`
}

func main() {
	const port = 8080
	http.HandleFunc("/plain", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "plain-tag")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Tag{Tag: "json-tag"})
		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Serving at port %d...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
