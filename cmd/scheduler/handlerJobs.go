package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

func (srv *server) handlerGetJobs(w http.ResponseWriter, req *http.Request) {
	jobs, err := srv.cfg.db.GetJobs(req.Context())
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't load Jobs from DB")
		return
	}

	type respData struct {
		Jobs []database.Job `json:"jobs"`
	}
	data := respData{
		Jobs: jobs,
	}
	requests.RespondWithJSON(w, http.StatusOK, data)
}

func (s *server) handlerDeleteJobs(w http.ResponseWriter, req *http.Request) {
	err := s.cfg.db.DeleteJobs(req.Context())
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

func (srv *server) handlerCreateJob(w http.ResponseWriter, req *http.Request) {

	requestData, err := requests.DecodeRequest(req, jobParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	job, err := srv.cfg.db.CreateJob(req.Context(), database.CreateJobParams{
		Name:     requestData.Name,
		Status:   int32(utils.StatusWaiting),
		Priority: int32(requestData.Priority),
	})
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't create Job")
		return
	}

	err = createTasks(srv, job.ID, requestData.Tasks, req)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't create Tasks")
		return
	}

	tasks, _ := srv.cfg.db.GetTasks(req.Context())
	requests.RespondWithJSON(w, http.StatusOK, tasks)
}

func createTasks(srv *server, jobID uuid.UUID, tasks []taskParams, req *http.Request) error {

	if len(tasks) == 0 || tasks == nil {
		return nil
	}

	for _, t := range tasks {

		_, err := srv.cfg.db.CreateTask(req.Context(), database.CreateTaskParams{
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
