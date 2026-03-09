package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type worker struct {
	id     uuid.UUID
	conn   *websocket.Conn
	status utils.ObjectStatus
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wrk := worker{}

	host := os.Getenv("TS_HOST")
	port := os.Getenv("TS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	ws, err := wrk.connectToWs(addr)
	if err != nil {
		log.Fatal(err)
	}
	wrk.conn = ws
	wg := sync.WaitGroup{}
	wg.Add(1)

	wrk.sendConnectMessage()
	go wrk.readWsMessages(&wg)

	wg.Wait()
}

func (wrk *worker) sendConnectMessage() {
	type connectMessage struct {
		Type    int               `json:"message_type"`
		Payload map[string]string `json:"payload"`
	}

	payload := connectMessage{
		Type:    int(utils.ConnectMessage),
		Payload: make(map[string]string, 0),
	}

	payload.Payload["host"] = "localhost"
	payload.Payload["port"] = "941241"

	message := requests.EncodeJSON(payload)

	wrk.conn.WriteMessage(websocket.BinaryMessage, message)
}

func (wrk *worker) readWsMessages(wg *sync.WaitGroup) {
	type serverMessage struct {
		Type    int               `json:"message_type"`
		Payload map[string]string `json:"payload"`
	}

	defer wrk.conn.Close()
	defer wg.Done()
	for {
		_, msg, err := wrk.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected connection close")
				return
			}
			log.Println(err)
		}

		jsonMsg, err := requests.DecodeJSON(msg, serverMessage{})
		if err != nil {
			log.Println("Could not decode worker message")
		}

		switch jsonMsg.Type {
		case int(utils.TaskMessage):
			wrk.handleTaskMessage(jsonMsg.Payload)
		default:
			log.Println("Unknown Message type:", jsonMsg.Type)
		}
	}
}

func (wrk *worker) connectToWs(addr string) (*websocket.Conn, error) {

	url := url.URL{Scheme: "ws", Host: addr, Path: "/api/registerWorkers"}
	ws, resp, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fmt.Println(resp.Status)
		return nil, err
	}
	return ws, nil
}

func (wrk *worker) handleTaskMessage(payload map[string]string) {

	jobID := uuid.MustParse(payload["job_id"])
	taskID := uuid.MustParse(payload["task_id"])
	cmd := payload["command"]

	err := wrk.sendStatusMessage(utils.StatusRunning, jobID, taskID, nil)
	if err != nil {
		log.Println("Could not send status update. Aborting task launch.")
		return
	}
	go wrk.runTask(jobID, taskID, cmd)

}

func (wrk *worker) sendStatusMessage(status utils.ObjectStatus, jobID, taskID uuid.UUID, output *string) error {
	type statusMessage struct {
		Type    int               `json:"message_type"`
		Payload map[string]string `json:"payload"`
	}

	payload := statusMessage{
		Type:    int(utils.StatusMessage),
		Payload: make(map[string]string, 0),
	}

	payload.Payload["status"] = strconv.Itoa(int(status))
	payload.Payload["job_id"] = jobID.String()
	payload.Payload["task_id"] = taskID.String()
	if output != nil {
		payload.Payload["output"] = *output
	}

	message := requests.EncodeJSON(payload)

	err := wrk.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (wrk *worker) runTask(jobID, taskID uuid.UUID, taskCmd string) {

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

	var output string

	log.Println("Running task", taskCmd)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr.String())
		output = stderr.String()
		wrk.sendStatusMessage(utils.StatusError, jobID, taskID, &output)
		return
	}

	fmt.Println("Finished task", taskCmd)
	output = stdout.String()
	wrk.sendStatusMessage(utils.StatusFinished, jobID, taskID, &output)
}
