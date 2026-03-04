package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
	"github.com/ju-vfx/task-scheduler/internal/requests"
)

func (s *server) handlerGetWorkers(w http.ResponseWriter, req *http.Request) {
	type respData struct {
		Workers []database.Worker `json:"workers"`
	}

	data := respData{
		Workers: s.cfg.workers,
	}
	requests.RespondWithJSON(w, http.StatusOK, data)
}

func (s *server) handlerRegisterWorker(w http.ResponseWriter, req *http.Request) {

	type workerParams struct {
		ID   *string `json:"id"`
		Host string  `json:"host"`
		Port string  `json:"port"`
	}
	reqParams, err := requests.DecodeRequest(req, workerParams{})
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	if reqParams.ID != nil {
		err = s.cfg.db.UpdateLastSeen(req.Context(), uuid.MustParse(*reqParams.ID))
		if err != nil {
			requests.RespondWithError(w, http.StatusInternalServerError, "Error logging in worker")
			return
		}
		updateWorkerSlice(s)
	} else {
		worker, err := s.cfg.db.CreateWorker(req.Context(), database.CreateWorkerParams{Host: reqParams.Host, Port: reqParams.Port})
		if err != nil {
			requests.RespondWithError(w, http.StatusInternalServerError, "Error adding worker")
			return
		}

		type respParams struct {
			ID string `json:"id"`
		}
		updateWorkerSlice(s)
		requests.RespondWithJSON(w, http.StatusOK, respParams{ID: worker.ID.String()})
	}
}

func (s *server) handlerDeleteWorkers(w http.ResponseWriter, req *http.Request) {
	err := s.cfg.db.DeleteWorkers(req.Context())
	if err != nil {
		requests.RespondWithError(w, http.StatusBadRequest, "Error deleting workers")
		return
	}
	updateWorkerSlice(s)
	requests.RespondWithJSON(w, http.StatusOK, "")
}
