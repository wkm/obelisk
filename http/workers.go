package main

type WorkerStatus uint

var (
	StatusOk      = 0
	StatusError   = 1
	StatusUnknown = 2
)

type Worker struct {
	id, host string
	status   WorkerStatus
}
