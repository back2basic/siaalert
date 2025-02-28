package cron

import (
	"sync"
)

type Job struct {
	ID      int
	Name    string
	Type    string
	Address string
	HostKey string
	V2      bool
}

type Worker struct {
	ID        int
	JobQueue  chan Job
	Waitgroup *sync.WaitGroup
}
