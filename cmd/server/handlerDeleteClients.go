package main

import (
	"fmt"
	"net/http"
)

func (s *server) handlerDeleteClients(w http.ResponseWriter, req *http.Request) {
	err := s.db.DeleteClients(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error deleting clients.", err)
	}
}
