package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func newServer(cfg *appConfig) (*server, error) {
	srv := server{
		cfg: cfg,
	}
	return &srv, nil
}

func (srv *server) Start() {

	host := os.Getenv("TS_HOST")
	port := os.Getenv("TS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Starting Scheduler Server on http://%s", addr)

	registerHandlers(srv)
	updateWorkerSlice(srv)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func registerHandlers(srv *server) {
	// http.HandleFunc("GET /", s.handlerRoot)
	http.HandleFunc("GET /api/workers", srv.handlerGetWorkers)
	http.HandleFunc("POST /api/workers", srv.handlerRegisterWorker)
	if platform := os.Getenv("TS_PLATFORM"); platform == "dev" {
		http.HandleFunc("DELETE /api/workers", srv.handlerDeleteWorkers)
	}
}

func updateWorkerSlice(srv *server) {
	workers, err := srv.cfg.db.GetWorkers(context.Background())
	if err != nil {
		log.Fatal("Could not get workers")
	}
	srv.cfg.workers = workers
}
