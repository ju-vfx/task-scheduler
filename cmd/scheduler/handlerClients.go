package main

import (
	"context"
	"log"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type client struct {
	conn *websocket.Conn
	conf *appConfig
	id   uuid.UUID
}

type clientMessage struct {
	Type    string      `json:"message_type"`
	Payload interface{} `json:"payload"`
}

func (conf *appConfig) handlerRegisterClients(w http.ResponseWriter, req *http.Request) {
	ws, err := UpgradeConnection(w, req)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Could not connect to websocket")
		return
	}

	client := &client{
		conn: ws,
		conf: conf,
		id:   uuid.New(),
	}
	conf.clients = append(conf.clients, client)

}

func (c *client) ReadClientWebsocketMessage() {
	defer c.conn.Close()
	for {
		messageType, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				c.conf.deleteClient(c)
				return
			}
			log.Println("Error reading message:", messageType, err)
			return
		}
	}
}

func (conf *appConfig) deleteClient(c *client) {
	conf.mu.Lock()
	conf.clients = slices.DeleteFunc(conf.clients, func(client *client) bool { return client.id == c.id })
	conf.mu.Unlock()
}

func (conf *appConfig) broadcastClients(payload interface{}) {
	message := requests.EncodeJSON(payload)
	for _, client := range conf.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
		}
	}
}

func (conf *appConfig) broadcastJobs() {
	dbJobs, err := conf.db.GetJobs(context.Background())
	if err != nil {
		log.Println("Can't load Jobs from DB", err)
		return
	}

	type respTask struct {
		TaskID      string `json:"task_id"`
		TaskName    string `json:"task_name"`
		TaskStatus  string `json:"task_status"`
		TaskCommand string `json:"task_command"`
		CreatedAt   string `json:"task_created_at"`
		FinishedAt  string `json:"task_finished_at"`
		Output      string `json:"task_output"`
	}
	type respJob struct {
		JobID         string     `json:"job_id"`
		JobName       string     `json:"job_name"`
		JobPriority   int        `json:"job_priority"`
		JobStatus     string     `json:"job_status"`
		JobCreatedAt  string     `json:"job_created_at"`
		JobFinishedAt string     `json:"job_finished_at"`
		JobTasks      []respTask `json:"job_tasks"`
	}
	data := make([]respJob, 0)
	for _, job := range dbJobs {
		dbTasks, err := conf.db.GetTasksByJobId(context.Background(), job.ID)
		if err != nil {
			log.Println("Can't load Tasks from DB", err)
			return
		}
		tasks := make([]respTask, 0)
		for _, task := range dbTasks {
			output := ""
			if task.Stdout.Valid {
				output = task.Stdout.String
			} else if task.Stderr.Valid {
				output = task.Stderr.String
			}
			t := respTask{
				TaskID:      task.ID.String(),
				TaskName:    task.Name,
				TaskStatus:  utils.ObjectStatus(task.Status).String(),
				TaskCommand: task.Command,
				CreatedAt:   utils.TimeToString(task.CreatedAt),
				FinishedAt:  utils.TimeToString(task.FinishedAt.Time),
				Output:      output,
			}

			tasks = append(tasks, t)
		}
		j := respJob{
			JobID:         job.ID.String(),
			JobName:       job.Name,
			JobPriority:   int(job.Priority),
			JobStatus:     utils.ObjectStatus(job.Status).String(),
			JobCreatedAt:  utils.TimeToString(job.CreatedAt),
			JobFinishedAt: utils.TimeToString(job.FinishedAt.Time),
			JobTasks:      tasks,
		}
		data = append(data, j)
	}

	conf.broadcastClients(clientMessage{Type: "jobs", Payload: data})
}

func (conf *appConfig) broadcastWorkers() {
	type workerType struct {
		ID          string `json:"id"`
		Host        string `json:"host"`
		Status      string `json:"status"`
		LastSeenAt  string `json:"last_seen_at"`
		ConnectedAt string `json:"connected_at"`
	}
	payload := make([]workerType, 0)
	for _, worker := range conf.workers {
		w := workerType{
			ID:          worker.id.String(),
			Host:        worker.host,
			Status:      worker.status.String(),
			LastSeenAt:  utils.TimeToString(worker.lastSeenAt),
			ConnectedAt: utils.TimeToString(worker.connectedAt),
		}

		payload = append(payload, w)
	}

	conf.broadcastClients(clientMessage{Type: "workers", Payload: payload})
}
