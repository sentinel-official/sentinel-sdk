package cron

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Scheduler manages the scheduling and execution of workers.
type Scheduler struct {
	isRunning  bool              // Indicates if the scheduler is currently running.
	stopSignal chan struct{}     // Channel to signal workers to stop.
	workers    map[string]Worker // Workers registered with the scheduler.
	mu         sync.Mutex        // Mutex for synchronizing access to scheduler state.
	wg         sync.WaitGroup    // WaitGroup to wait for all worker goroutines to complete.
}

// NewScheduler creates and initializes a new Scheduler instance.
func NewScheduler() *Scheduler {
	return &Scheduler{
		stopSignal: make(chan struct{}),
		workers:    make(map[string]Worker),
	}
}

// Start begins executing all registered workers in separate goroutines.
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return errors.New("scheduler is already running")
	}

	s.isRunning = true

	// Run each worker in its own goroutine.
	for _, w := range s.workers {
		s.wg.Add(1)
		go s.runWorker(w)
	}

	return nil
}

// Stop halts the execution of all workers and stops the scheduler.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	// Wait for all worker goroutines to finish.
	close(s.stopSignal)
	s.wg.Wait()

	s.isRunning = false
}

// RegisterWorkers adds multiple workers to the scheduler.
func (s *Scheduler) RegisterWorkers(workers ...Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return errors.New("scheduler is already running")
	}

	for _, w := range workers {
		if _, exists := s.workers[w.Name()]; exists {
			return fmt.Errorf("duplicate worker %s", w.Name())
		}

		s.workers[w.Name()] = w
	}

	return nil
}

// runWorker continuously executes a worker's function and handles errors.
func (s *Scheduler) runWorker(w Worker) {
	defer func() {
		w.OnExit()
		s.wg.Done()
	}()

	for retries := 0; ; {
		select {
		case <-s.stopSignal:
			return
		default:
			if err := w.Run(); err != nil {
				if stop := w.OnError(err); stop {
					return
				}

				retries++
				if retries == w.MaxRetries() {
					return
				}
			} else {
				retries = 0

				// Sleep before the next execution.
				time.Sleep(w.Interval())
			}
		}
	}
}
