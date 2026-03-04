package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
)

func (s *server) handlerGetJobs(w http.ResponseWriter, req *http.Request) {
	type respData struct {
		Workers []database.Worker `json:"workers"`
	}

	data := respData{
		Workers: s.cfg.workers,
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
	Name       string        `json:"name"`
	Command    string        `json:"command"`
	ChildTasks *[]taskParams `json:"child_tasks"`
}
type jobParams struct {
	Name     string       `json:"name"`
	Priority int          `json:"priority"`
	Tasks    []taskParams `json:"tasks"`
}

func (s *server) handlerCreateJob(w http.ResponseWriter, req *http.Request) {

	requestData, err := requests.DecodeRequest(req, jobParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	job, err := s.cfg.db.CreateJob(req.Context(), database.CreateJobParams{
		Name:     requestData.Name,
		Status:   "waiting",
		Priority: int32(requestData.Priority),
	})
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't create Job")
		return
	}

	err = createTasksRecursive(s, job.ID, uuid.NullUUID{}, requestData.Tasks, req)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't create Tasks")
		return
	}

	tasks, _ := s.cfg.db.GetTasks(req.Context())
	fmt.Println("Number of tasks: ", len(tasks))
	requests.RespondWithJSON(w, http.StatusOK, tasks)
}

func createTasksRecursive(s *server, jobID uuid.UUID, parentTaskID uuid.NullUUID, tasks []taskParams, req *http.Request) error {

	if len(tasks) == 0 || tasks == nil {
		return nil
	}

	for _, t := range tasks {

		task, err := s.cfg.db.CreateTask(req.Context(), database.CreateTaskParams{
			Name:         t.Name,
			Status:       "waiting",
			ParentTaskID: parentTaskID,
			Command:      t.Command,
			JobID:        jobID,
		})
		if err != nil {
			return err
		}
		fmt.Println(t.Name, t.ChildTasks)
		if t.ChildTasks != nil {
			err = createTasksRecursive(s, jobID, uuid.NullUUID{UUID: task.ID, Valid: true}, *t.ChildTasks, req)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
