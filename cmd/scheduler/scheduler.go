package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type job struct {
	job   database.Job
	tasks []database.Task
}

func (conf *appConfig) ScheduleTasks() {
	conf.mu.Lock()
	defer conf.mu.Unlock()

	availableWorkers := getAvailableWorkers(conf.workers)
	if len(availableWorkers) < 1 {
		return
	}

	waitingJobs := getWaitingJobs(conf.db)
	if len(waitingJobs) < 1 {
		return
	}

	runNextTask(availableWorkers, waitingJobs)
}

func getAvailableWorkers(workers []*worker) []*worker {
	availableWorkers := make([]*worker, 0)
	for _, worker := range workers {
		if worker.status == utils.StatusWaiting {
			availableWorkers = append(availableWorkers, worker)
		}
	}
	return availableWorkers
}

func getWaitingJobs(db *database.Queries) []job {
	jobs := make([]job, 0)
	waitingJobs, err := db.GetWaitingJobs(context.Background())
	if err != nil {
		log.Fatal("Can't load waiting jobs")
	}
	for _, waitingJob := range waitingJobs {
		tasks, err := db.GetTasksByJobId(context.Background(), waitingJob.ID)
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

func runNextTask(availableWorkers []*worker, waitingJobs []job) {
workerLoop:
	for _, worker := range availableWorkers {
		for _, job := range waitingJobs {
			for _, task := range job.tasks {
				if task.Status == int32(utils.StatusWaiting) {
					err := sendTaskToWorker(worker, job.job.ID, task.ID, task.Command)
					if err != nil {
						log.Printf("Can't send task %s to worker.", task.ID)
						continue
					}
					worker.task_id = &task.ID
					task.Status = int32(utils.StatusRunning)
					log.Println("Sending task", task.Name, "to", worker.host)
					continue workerLoop
				}
			}
		}
	}
}

func sendTaskToWorker(w *worker, jobID, taskID uuid.UUID, cmd string) error {
	err := w.sendTaskMessage(jobID, taskID, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (conf *appConfig) updateJobStatus(jobID uuid.UUID) {

	tasks, err := conf.db.GetTasksByJobId(context.Background(), jobID)
	if err != nil {
		log.Println(err)
	}

	jobStatus := calcJobStatus(tasks)

	updateParams := database.UpdateJobStatusParams{ID: jobID}

	switch jobStatus {
	case utils.StatusFinished:
		updateParams.Status = int32(jobStatus)
		updateParams.FinishedAt.Time = time.Now()
		updateParams.FinishedAt.Valid = true
	case utils.StatusError:
		updateParams.Status = int32(jobStatus)
	default:
		updateParams.Status = int32(utils.StatusRunning)
	}

	err = conf.db.UpdateJobStatus(context.Background(), updateParams)
	if err != nil {
		log.Println(err)
	}
}

func calcJobStatus(tasks []database.Task) utils.ObjectStatus {
	waitingCount := 0
	runningCount := 0
	finishedCount := 0
	errorCount := 0

	for _, task := range tasks {
		switch utils.ObjectStatus(task.Status) {
		case utils.StatusWaiting:
			waitingCount++
		case utils.StatusRunning:
			runningCount++
		case utils.StatusFinished:
			finishedCount++
		case utils.StatusError:
			errorCount++
		}
	}

	jobStatus := utils.StatusWaiting
	if errorCount > 0 {
		jobStatus = utils.StatusError
	} else if finishedCount == len(tasks) {
		jobStatus = utils.StatusFinished
	} else if runningCount > 0 && waitingCount > 0 {
		jobStatus = utils.StatusRunning
	}

	return jobStatus
}

func (conf *appConfig) updateTaskStatus(taskID uuid.UUID, status utils.ObjectStatus, output string) {

	taskStatusParms := database.UpdateTaskStatusParams{ID: taskID}
	switch status {
	case utils.StatusFinished:
		taskStatusParms.Status = int32(status)
		taskStatusParms.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
		taskStatusParms.Stdout = sql.NullString{String: output, Valid: true}
	default:
		taskStatusParms.Status = int32(status)
		taskStatusParms.CancelledAt = sql.NullTime{Time: time.Now(), Valid: true}
		taskStatusParms.Stderr = sql.NullString{String: output, Valid: true}
	}
	_, err := conf.db.UpdateTaskStatus(context.Background(), taskStatusParms)
	if err != nil {
		log.Println("Can't update task status.", err)
		return
	}
}
