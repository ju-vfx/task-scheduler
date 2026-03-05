package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ju-vfx/task-scheduler/internal/database"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}

	conf := &appConfig{
		db: db,
	}

	server, err := newServer(conf)
	if err != nil {
		log.Fatal(err)
	}

	scheduler, err := newScheduler(conf)
	if err != nil {
		log.Fatal(err)
	}

	go server.Start()
	scheduler.Start()
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
