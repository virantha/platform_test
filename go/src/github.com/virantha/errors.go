package main

import (
	"log"
	"encoding/json"
	"net/http"
)

// Write out the proper json error message
func WriteJSONError(w http.ResponseWriter, errorCode int, message string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	j, err := json.Marshal(struct {
		Error string `json:"error"`
	}{message})
	if err != nil {
		// Hmm, I'd better not have an error here!
		log.Println(err)
	}

	w.Write(j)
}
