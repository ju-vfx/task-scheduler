package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ju-vfx/task-scheduler/internal/database"
)

func (s *server) handlerGetWorkers(w http.ResponseWriter, req *http.Request) {
	type respData struct {
		Workers []database.Worker `json:"workers"`
	}

	data := respData{
		Workers: s.workers,
	}
	respondWithJSON(w, http.StatusOK, data)
}

func (s *server) handlerRegisterWorker(w http.ResponseWriter, req *http.Request) {

	type workerParams struct {
		ID   *string `json:"id"`
		Host string  `json:"host"`
		Port string  `json:"port"`
	}
	reqParams, err := decodeRequest(req, workerParams{})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't decode Request Body")
		return
	}

	if reqParams.ID != nil {
		err = s.db.UpdateLastSeen(req.Context(), uuid.MustParse(*reqParams.ID))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error logging in worker")
			return
		}
		updateWorkerSlice(s)
	} else {
		worker, err := s.db.CreateWorker(req.Context(), database.CreateWorkerParams{Host: reqParams.Host, Port: reqParams.Port})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error adding worker")
			return
		}

		type respParams struct {
			ID string `json:"id"`
		}
		updateWorkerSlice(s)
		respondWithJSON(w, http.StatusOK, respParams{ID: worker.ID.String()})
	}
}

func (s *server) handlerDeleteWorkers(w http.ResponseWriter, req *http.Request) {
	err := s.db.DeleteWorkers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error deleting workers")
		return
	}
	updateWorkerSlice(s)
	respondWithJSON(w, http.StatusOK, "")
}

func decodeRequest[T interface{}](req *http.Request, i T) (T, error) {
	params := i
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}
