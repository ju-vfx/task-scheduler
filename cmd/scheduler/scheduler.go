package main

import (
	"bytes"
	"context"
	"database/sql"
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
		sdl.updateJobStatus(waitingJobs)

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
						fmt.Println("Error sending job to worker:", err)
						break jobLoop
					}
					log.Printf("Sending task %s to worker %s\n", task.Name, worker.Host)
					err = sdl.cfg.db.UpdateWorkerStatus(context.Background(), database.UpdateWorkerStatusParams{ID: worker.ID, Status: int32(utils.StatusRunning)})
					if err != nil {
						fmt.Println(err)
					}
					_, err = sdl.cfg.db.UpdateTaskStatus(context.Background(), database.UpdateTaskStatusParams{ID: task.ID, Status: int32(utils.StatusRunning)})
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

func (sdl *scheduler) updateJobStatus(waitingJobs []job) {
	for _, job := range waitingJobs {

		jobStatus := utils.ObjectStatus(job.job.Status)
		finishedCount := 0
		runningCount := 0
		waitingCount := 0
		errorCount := 0

		for _, task := range job.tasks {
			switch utils.ObjectStatus(task.Status) {
			case utils.StatusRunning:
				runningCount++
			case utils.StatusWaiting:
				waitingCount++
			case utils.StatusFinished:
				finishedCount++
			case utils.StatusError:
				errorCount++
			default:
				errorCount++
			}
		}

		// fmt.Println("Running:", runningCount, "Waiting:", waitingCount, "Finished:", finishedCount, "Error:", errorCount)
		updateParams := database.UpdateJobStatusParams{ID: job.job.ID}
		if finishedCount == len(job.tasks) {
			jobStatus = utils.StatusFinished
			updateParams.Status = int32(jobStatus)
			updateParams.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
		} else if runningCount > 0 && errorCount == 0 {
			jobStatus = utils.StatusRunning
			updateParams.Status = int32(jobStatus)
		} else if runningCount == 0 && errorCount == 0 {
			jobStatus = utils.StatusWaiting
			updateParams.Status = int32(jobStatus)
		} else if errorCount > 0 {
			jobStatus = utils.StatusError
			updateParams.Status = int32(jobStatus)
			updateParams.CancelledAt = sql.NullTime{Time: time.Now(), Valid: true}
		}

		err := sdl.cfg.db.UpdateJobStatus(context.Background(), updateParams)
		if err != nil {
			log.Println("Could not update Job status:", err)
		}
	}

}
