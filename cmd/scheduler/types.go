package main

import "github.com/ju-vfx/task-scheduler/internal/database"

type appConfig struct {
	db      *database.Queries
	workers []database.Worker
	jobs    []job
}

type server struct {
	cfg *appConfig
}

type scheduler struct {
	cfg *appConfig
}

type job struct {
	job   database.Job
	tasks []database.Task
}
