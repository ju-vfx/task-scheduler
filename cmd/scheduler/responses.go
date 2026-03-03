package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type error struct {
		Error string `json:"error"`
	}

	respBody := error{
		Error: msg,
	}

	respData, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error Mashalling JSON for Error\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respData)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	respData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Mashalling JSON for Response\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respData)
}
