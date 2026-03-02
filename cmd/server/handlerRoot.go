package main

import (
	"io"
	"net/http"
)

func (s *server) handlerRoot(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "/")
}
