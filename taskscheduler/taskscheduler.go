// Package taskscheduler is a package for a distributed task scheduler
package taskscheduler

// TODO implement distributed task scheduling

import (
	"sync"
)

type Server struct {
	// port
	// some way to gather metrics (cpu, memory, liveliness)
	//
}

type Job interface {
	Run()
}

type Batch struct {
	Jobs  []Job
	Count int
}

type TaskScheduler struct {
	wg      sync.WaitGroup
	servers []Server
}

func (ts *TaskScheduler) Submit(batch *Batch) error {
	// novel implementaiton
	// TODO implement actual scheduling
	for jobIdx := range batch.Count {
		ts.wg.Go(batch.Jobs[jobIdx].Run)
	}

	return nil
}

func (ts *TaskScheduler) Close() {
	ts.wg.Wait()
}
