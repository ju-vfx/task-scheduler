package main

import (
	"fmt"
	"net/http"

	"github.com/ju-vfx/task-scheduler/internal/database"
)

func (s *server) handlerAddClient(w http.ResponseWriter, req *http.Request) {
	host := "localhost"
	ipAddr := "192.168.178.3"
	client, err := s.db.CreateClient(req.Context(), database.CreateClientParams{Host: host, IpAddr: ipAddr})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error adding client", err)
	}

	fmt.Println(client)
}
