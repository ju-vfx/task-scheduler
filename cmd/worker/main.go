package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type worker struct {
	id      uuid.UUID
	host    string
	port    string
	srvHost string
	srvPort string
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wrk := worker{}

	wrk.srvHost = os.Getenv("TS_HOST")
	wrk.srvPort = os.Getenv("TS_PORT")

	wrk.host = os.Getenv("TS_WORKER_HOST")
	wrk.port = os.Getenv("TS_WORKER_PORT")
	addr := fmt.Sprintf("%s:%s", wrk.host, wrk.port)

	wrk.registerWithServer()

	http.HandleFunc("POST /api/tasks", wrk.handlerLaunchTask)

	log.Printf("Starting Worker on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (wrk *worker) handlerLaunchTask(w http.ResponseWriter, req *http.Request) {
	type launchTaskReq struct {
		ID      string `json:"id"`
		Command string `json:"command"`
	}
	params, err := requests.DecodeRequest(req, launchTaskReq{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Could not decode request")
		return
	}

	requests.RespondWithJSON(w, http.StatusOK, params)
	go wrk.runTask(uuid.MustParse(params.ID), params.Command)
}

func (wrk *worker) runTask(taskID uuid.UUID, taskCmd string) {
	fmt.Println("Running task ", taskCmd)
	time.Sleep(time.Second * 10) // simulate running command

	wrk.updateTaskStatus(taskID, utils.StatusFinished)
}

func (wrk *worker) updateTaskStatus(taskId uuid.UUID, status utils.ObjectStatus) {
	type loginParams struct {
		ID     string `json:"id"`
		TaskID string `json:"task_id"`
		Status int32  `json:"status"`
	}

	data, _ := json.Marshal(loginParams{
		ID:     wrk.id.String(),
		TaskID: taskId.String(),
		Status: int32(status),
	})

	_, err := http.Post(fmt.Sprintf("http://%s:%s/api/tasks", wrk.srvHost, wrk.srvPort), "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrk *worker) registerWithServer() {

	type loginParams struct {
		ID   *string `json:"id"`
		Host string  `json:"host"`
		Port string  `json:"port"`
	}

	data, _ := json.Marshal(loginParams{
		ID:   nil,
		Host: wrk.host,
		Port: wrk.port,
	})

	resp, err := http.Post(fmt.Sprintf("http://%s:%s/api/workers", wrk.srvHost, wrk.srvPort), "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	type respParms struct {
		ID string `json:"id"`
	}
	var parms respParms
	err = json.NewDecoder(resp.Body).Decode(&parms)
	if err != nil {
		log.Fatal(err)
	}

	wrk.id = uuid.MustParse(parms.ID)
}
