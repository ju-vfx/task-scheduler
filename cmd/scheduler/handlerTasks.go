package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

func (srv *server) handlerUpdateTasks(w http.ResponseWriter, req *http.Request) {
	type updateTaskParams struct {
		ID     string `json:"id"`
		TaskID string `json:"task_id"`
		Status int32  `json:"status"`
	}

	requestData, err := requests.DecodeRequest(req, updateTaskParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}
	status := utils.ObjectStatus(requestData.Status)
	workerStatus := status
	if status != utils.StatusError {
		workerStatus = utils.StatusWaiting
	}
	err = srv.cfg.db.UpdateWorkerStatus(req.Context(), database.UpdateWorkerStatusParams{ID: uuid.MustParse(requestData.ID), Status: int32(workerStatus)})
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't update worker status")
		return
	}
	err = srv.cfg.db.UpdateTaskStatus(req.Context(), database.UpdateTaskStatusParams{ID: uuid.MustParse(requestData.TaskID), Status: int32(status)})
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't update task status")
		return
	}
}
