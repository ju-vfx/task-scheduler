package main

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

func (conf *appConfig) handlerDeleteJobs(w http.ResponseWriter, req *http.Request) {
	err := conf.db.DeleteJobs(req.Context())
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't delete jobs")
		return
	}
}

type taskParams struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}
type jobParams struct {
	Name     string       `json:"name"`
	Priority int          `json:"priority"`
	Tasks    []taskParams `json:"tasks"`
}

func (conf *appConfig) handlerCreateJob(w http.ResponseWriter, req *http.Request) {

	buff, _ := io.ReadAll(req.Body)
	requestData, err := requests.DecodeJSON(buff, jobParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	job, err := conf.db.CreateJob(req.Context(), database.CreateJobParams{
		Name:     requestData.Name,
		Status:   int32(utils.StatusWaiting),
		Priority: int32(requestData.Priority),
	})
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't create Job")
		return
	}

	err = conf.createTasks(job.ID, requestData.Tasks, req)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't create Tasks")
		return
	}

	requests.RespondWithJSON(w, http.StatusOK, job)
	conf.broadcastJobs()
	conf.ScheduleTasks()
}

func (conf *appConfig) createTasks(jobID uuid.UUID, tasks []taskParams, req *http.Request) error {

	if len(tasks) == 0 || tasks == nil {
		return nil
	}

	for _, t := range tasks {

		_, err := conf.db.CreateTask(req.Context(), database.CreateTaskParams{
			Name:    t.Name,
			Status:  int32(utils.StatusWaiting),
			Command: t.Command,
			JobID:   jobID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
