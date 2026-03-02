package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type server struct {
	conf string
	db   string
}

func newServer() *server {
	srv := server{
		conf: "testconf",
		db:   "testdb",
	}
	return &srv
}

func (s *server) Start() {

	host := os.Getenv("TS_HOST")
	port := os.Getenv("TS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Starting Server on http://%s", addr)

	registerHandlers(s)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func registerHandlers(s *server) {
	http.HandleFunc("GET /", s.handlerRoot)
}
