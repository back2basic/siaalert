// scheduler/scheduler.go
package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a new instance of Scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
	}
}

// AddJob adds a new cron job with the specified schedule and job function
func (s *Scheduler) AddJob(schedule string, job func()) (cron.EntryID, error) {
	id, err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Start starts the cron scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Scheduler started")
}

// Stop stops the cron scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}
