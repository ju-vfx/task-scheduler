package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

func newScheduler(cfg *appConfig) (*scheduler, error) {
	sdl := scheduler{
		cfg:                   cfg,
		updateIntervalSeconds: 2,
	}
	return &sdl, nil
}

func (sdl *scheduler) Start() {

	scheduleLoop(sdl, time.Second*time.Duration(sdl.updateIntervalSeconds))
}

func scheduleLoop(sdl *scheduler, refreshInterval time.Duration) {
	for {
		time.Sleep(refreshInterval)
		waitingJobs := getWaitingJobs(sdl)
		if len(waitingJobs) < 1 {
			continue
		}
		availableWorkers := getAvailableWorkers(sdl)
		if len(availableWorkers) < 1 {
			continue
		}

		distributeTasks(sdl, waitingJobs, availableWorkers)
	}
}

func getWaitingJobs(sdl *scheduler) []job {
	jobs := make([]job, 0)
	waitingJobs, err := sdl.cfg.db.GetWaitingJobs(context.Background())
	if err != nil {
		log.Fatal("Can't load waiting jobs")
	}
	for _, waitingJob := range waitingJobs {
		tasks, err := sdl.cfg.db.GetTasksByJobId(context.Background(), waitingJob.ID)
		if err != nil {
			log.Fatal("Error getting Tasks for Job")
		}
		j := job{
			job:   waitingJob,
			tasks: tasks,
		}

		jobs = append(jobs, j)
	}

	return jobs
}

func getAvailableWorkers(sdl *scheduler) []database.Worker {
	availableWorkers := make([]database.Worker, 0)
	workers, err := sdl.cfg.db.GetWorkers(context.Background())
	if err != nil {
		log.Fatal("Can't get workers from DB")
	}
	for _, worker := range workers {
		if worker.Status == int32(utils.StatusWaiting) {
			availableWorkers = append(availableWorkers, worker)
		}
	}
	return availableWorkers
}

func distributeTasks(sdl *scheduler, waitingJobs []job, availableWorkers []database.Worker) {
	for _, worker := range availableWorkers {
	jobLoop:
		for jobIdx := range waitingJobs {
			for taskIdx, task := range waitingJobs[jobIdx].tasks {
				if task.Status == int32(utils.StatusWaiting) {
					waitingJobs[jobIdx].tasks[taskIdx].Status = int32(utils.StatusRunning)
					err := sdl.sendTaskToWorker(task, worker)
					if err != nil {
						break jobLoop
					}
					log.Printf("Sending task %s to worker %s\n", task.Name, worker.Host)
					err = sdl.cfg.db.UpdateWorkerStatus(context.Background(), database.UpdateWorkerStatusParams{ID: worker.ID, Status: int32(utils.StatusRunning)})
					if err != nil {
						fmt.Println(err)
					}
					err = sdl.cfg.db.UpdateTaskStatus(context.Background(), database.UpdateTaskStatusParams{ID: task.ID, Status: int32(utils.StatusRunning)})
					if err != nil {
						fmt.Println(err)
					}
					break jobLoop
				}
			}
		}
	}
}

func (sdl *scheduler) sendTaskToWorker(task database.Task, worker database.Worker) error {

	type loginParams struct {
		ID      string `json:"id"`
		Command string `json:"command"`
	}

	data, _ := json.Marshal(loginParams{
		ID:      task.ID.String(),
		Command: task.Command,
	})

	resp, err := http.Post(fmt.Sprintf("http://%s:%s/api/tasks", worker.Host, worker.Port), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Task not sent successfully")
	}

	return nil
}
