package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ju-vfx/task-scheduler/internal/database"
	_ "github.com/lib/pq"
)

type server struct {
	conf string
	db   *database.Queries
}

func newServer() (*server, error) {
	srv := server{
		conf: "testconf",
		db:   connectDb(),
	}
	return &srv, nil
}

func (s *server) Start() {

	host := os.Getenv("TS_HOST")
	port := os.Getenv("TS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Starting Server on http://%s", addr)

	registerHandlers(s)

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
	http.HandleFunc("GET /", s.handlerRoot)
}
