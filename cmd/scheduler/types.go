package main

import "github.com/ju-vfx/task-scheduler/internal/database"

type server struct {
	db      *database.Queries
	workers []database.Worker
}
