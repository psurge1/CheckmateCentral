// Package taskscheduler is a package for a distributed task scheduler
package taskscheduler

// TODO implement distributed task scheduling

import (
	"sync"
)

type Job interface {
	Run()
}

type TaskScheduler struct {
	wg sync.WaitGroup
}

func (ts *TaskScheduler) Submit(job Job) {
	ts.wg.Go(job.Run)
}

func (ts *TaskScheduler) Close() {
	ts.wg.Wait()
}
