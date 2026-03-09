package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/ju-vfx/task-scheduler/internal/database"
	_ "github.com/lib/pq"
)

type appConfig struct {
	mu      sync.Mutex
	db      *database.Queries
	workers []*worker
	jobs    []*job
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}

	conf := appConfig{
		db:      db,
		workers: make([]*worker, 0),
		jobs:    make([]*job, 0),
	}

	conf.startServer()
}

func (c *appConfig) TryMe() {

}

func connectDb() (*database.Queries, error) {
	host := os.Getenv("TS_DB_HOST")
	port := os.Getenv("TS_DB_PORT")
	user := os.Getenv("TS_DB_USER")
	password := os.Getenv("TS_DB_PASSWORD")
	dbName := os.Getenv("TS_DB_NAME")
	sslMode := os.Getenv("TS_DB_SSLMODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return database.New(db), nil
}

func (conf *appConfig) startServer() {

	host := os.Getenv("TS_HOST")
	port := os.Getenv("TS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Starting Scheduler Server on http://%s", addr)

	conf.registerHandlers()
	if platform := os.Getenv("TS_PLATFORM"); platform == "dev" {
		_ = conf.db.DeleteJobs(context.Background())
		_ = conf.db.DeleteWorkers(context.Background())
	}

	log.Fatal(http.ListenAndServe(addr, nil))
}

func (conf *appConfig) registerHandlers() {
	http.HandleFunc("GET /api/workers", conf.handlerGetWorkers)
	http.HandleFunc("/api/registerWorkers", conf.handlerRegisterWorker)

	http.HandleFunc("GET /api/jobs", conf.handlerGetJobs)
	http.HandleFunc("POST /api/jobs", conf.handlerCreateJob)

	http.HandleFunc("POST /api/tasks", conf.handlerUpdateTasks)

	if platform := os.Getenv("TS_PLATFORM"); platform == "dev" {
		http.HandleFunc("DELETE /api/jobs", conf.handlerDeleteJobs)
	}
}
