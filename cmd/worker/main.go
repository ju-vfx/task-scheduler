package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("TS_WORKER_HOST")
	port := os.Getenv("TS_WORKER_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	registerWithServer(host, port)

	log.Printf("Starting Worker on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func registerWithServer(host, port string) {

	type loginParams struct {
		ID   *string `json:"id"`
		Host string  `json:"host"`
		Port string  `json:"port"`
	}

	data, _ := json.Marshal(loginParams{
		ID:   nil,
		Host: host,
		Port: port,
	})

	http.Post(fmt.Sprintf("http://%s:%s/api/workers", os.Getenv("TS_HOST"), os.Getenv("TS_PORT")), "application/json", bytes.NewReader(data))
}
