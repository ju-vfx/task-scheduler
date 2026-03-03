package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ju-vfx/task-scheduler/internal/database"
	_ "github.com/lib/pq"
)

func newServer() (*server, error) {
	srv := server{
		db:      connectDb(),
		workers: make([]database.Worker, 0),
	}
	return &srv, nil
}

func (s *server) Start() {

	host := os.Getenv("TS_HOST")
	port := os.Getenv("TS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Starting Scheduler Server on http://%s", addr)

	registerHandlers(s)
	updateWorkerSlice(s)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func connectDb() *database.Queries {
	host := os.Getenv("TS_DB_HOST")
	port := os.Getenv("TS_DB_PORT")
	user := os.Getenv("TS_DB_USER")
	password := os.Getenv("TS_DB_PASSWORD")
	dbName := os.Getenv("TS_DB_NAME")
	sslMode := os.Getenv("TS_DB_SSLMODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	return database.New(db)
}

func registerHandlers(s *server) {
	// http.HandleFunc("GET /", s.handlerRoot)
	http.HandleFunc("GET /api/workers", s.handlerGetWorkers)
	http.HandleFunc("POST /api/workers", s.handlerRegisterWorker)
	if platform := os.Getenv("TS_PLATFORM"); platform == "dev" {
		http.HandleFunc("DELETE /api/workers", s.handlerDeleteWorkers)
	}
}

func updateWorkerSlice(s *server) {
	workers, err := s.db.GetWorkers(context.Background())
	if err != nil {
		log.Fatal("Could not get workers")
	}
	s.workers = workers
}
