package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

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
	cmdSlice := strings.Fields(taskCmd)

	var cmd *exec.Cmd

	if len(cmdSlice) < 2 {
		cmd = exec.Command(cmdSlice[0])
	} else {
		cmd = exec.Command(cmdSlice[0], cmdSlice[1:]...)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Println("Running task", taskCmd)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr.String())
		wrk.updateTaskStatus(taskID, utils.StatusError, stderr)
		return
	}

	fmt.Println("Finished task", taskCmd)
	wrk.updateTaskStatus(taskID, utils.StatusFinished, stdout)
}

func (wrk *worker) updateTaskStatus(taskId uuid.UUID, status utils.ObjectStatus, output bytes.Buffer) {
	type loginParams struct {
		ID     string `json:"id"`
		TaskID string `json:"task_id"`
		Status int32  `json:"status"`
		Output string `json:"output"`
	}

	data, _ := json.Marshal(loginParams{
		ID:     wrk.id.String(),
		TaskID: taskId.String(),
		Status: int32(status),
		Output: output.String(),
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
