package main

import (
	"database/sql"
	"net/http"
	"time"

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
		Output string `json:"output"`
	}

	requestData, err := requests.DecodeRequest(req, updateTaskParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	status := utils.ObjectStatus(requestData.Status)

	err = srv.cfg.db.UpdateWorkerStatus(req.Context(), database.UpdateWorkerStatusParams{ID: uuid.MustParse(requestData.ID), Status: int32(utils.StatusWaiting)})
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't update worker status")
		return
	}
	taskStatusParms := database.UpdateTaskStatusParams{ID: uuid.MustParse(requestData.TaskID)}
	switch status {
	case utils.StatusFinished:
		taskStatusParms.Status = int32(status)
		taskStatusParms.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
		taskStatusParms.Stdout = sql.NullString{String: requestData.Output, Valid: true}
	default:
		taskStatusParms.Status = int32(status)
		taskStatusParms.CancelledAt = sql.NullTime{Time: time.Now(), Valid: true}
		taskStatusParms.Stderr = sql.NullString{String: requestData.Output, Valid: true}
	}
	_, err = srv.cfg.db.UpdateTaskStatus(req.Context(), taskStatusParms)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't update task status")
		return
	}
}
