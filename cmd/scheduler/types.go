package main

import "github.com/ju-vfx/task-scheduler/internal/database"

type server struct {
	cfg *appConfig
}

type scheduler struct {
	cfg                   *appConfig
	updateIntervalSeconds int
}

type job struct {
	job   database.Job
	tasks []database.Task
}
