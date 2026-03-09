package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type job struct {
	job   database.Job
	tasks []database.Task
}

func (conf *appConfig) UpdateState() {
	conf.mu.Lock()
	defer conf.mu.Unlock()

	availableWorkers := getAvailableWorkers(conf.workers)
	if len(availableWorkers) < 1 {
		return
	}

	sendTaskToWorker(availableWorkers[0])

	waitingJobs := getWaitingJobs(conf.db)
	if len(waitingJobs) < 1 {
		return
	}

}

func sendTaskToWorker(w *worker) {
	task_id := uuid.New()
	err := w.sendTaskMessage(task_id, "testcommand")
	if err != nil {
		log.Println(err)
	}
	w.task_id = &task_id
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

func getAvailableWorkers(workers []*worker) []*worker {
	availableWorkers := make([]*worker, 0)
	for _, worker := range workers {
		if worker.status == utils.StatusWaiting {
			availableWorkers = append(availableWorkers, worker)
		}
	}
	return availableWorkers
}
