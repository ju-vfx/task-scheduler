package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

func (conf *appConfig) handlerUpdateTasks(w http.ResponseWriter, req *http.Request) {
	type updateTaskParams struct {
		ID     string `json:"id"`
		TaskID string `json:"task_id"`
		Status int32  `json:"status"`
		Output string `json:"output"`
	}
	var buff bytes.Buffer
	_ = req.Write(&buff)
	requestData, err := requests.DecodeJSON(buff.Bytes(), updateTaskParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	status := utils.ObjectStatus(requestData.Status)

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
	_, err = conf.db.UpdateTaskStatus(req.Context(), taskStatusParms)
	if err != nil {
		requests.RespondWithError(w, http.StatusInternalServerError, "Can't update task status")
		return
	}
}
