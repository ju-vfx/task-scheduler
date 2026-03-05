package utils

type ObjectStatus int

const (
	StatusUnknown ObjectStatus = iota
	StatusWaiting
	StatusRunning
	StatusError
	StatusFinished
)
