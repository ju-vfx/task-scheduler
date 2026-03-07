package utils

import (
	"time"
)

type ObjectStatus int

const (
	StatusUnknown ObjectStatus = iota
	StatusWaiting
	StatusRunning
	StatusError
	StatusFinished
)

func (e ObjectStatus) String() string {
	switch e {
	case StatusUnknown:
		return "Unknown"
	case StatusWaiting:
		return "Waiting"
	case StatusRunning:
		return "Running"
	case StatusError:
		return "Error"
	case StatusFinished:
		return "Finished"
	default:
		return "Undefined"
	}
}

func TimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	layout := "15:04:05 02/01/2006"
	return t.Format(layout)
}
