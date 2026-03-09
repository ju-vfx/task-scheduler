package requests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func DecodeJSON[T interface{}](payload []byte, i T) (T, error) {
	params := i
	err := json.NewDecoder(bytes.NewReader(payload)).Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}

func EncodeJSON(payload interface{}) []byte {
	respData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Encoding JSON\n")
		return nil
	}
	return respData
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(respData)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	respData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Mashalling JSON for Response\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(respData)
}
