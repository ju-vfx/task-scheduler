package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

func (conf *appConfig) handlerGetJobs(w http.ResponseWriter, req *http.Request) {
	dbJobs, err := conf.db.GetJobs(req.Context())
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't load Jobs from DB")
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
		dbTasks, err := conf.db.GetTasksByJobId(req.Context(), job.ID)
		if err != nil {
			requests.RespondWithError(w, http.StatusInternalServerError, "Can't load Tasks from DB")
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
	requests.RespondWithJSON(w, http.StatusOK, data)
}

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

	requestData, err := requests.DecodeRequest(req, jobParams{})
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

	tasks, _ := conf.db.GetTasks(req.Context())
	requests.RespondWithJSON(w, http.StatusOK, tasks)
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
