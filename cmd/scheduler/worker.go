package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ju-vfx/task-scheduler/internal/utils"
)

type worker struct {
	id          uuid.UUID
	conn        *websocket.Conn
	host        string
	port        string
	connectedAt time.Time
	lastSeenAt  time.Time
	status      utils.ObjectStatus
	task        string
}
