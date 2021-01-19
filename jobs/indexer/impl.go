package indexer

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

type jobsIndexer struct {
	storage          storage.JobsStorage
	runContext       context.Context
	cancelRunContext func()
	done             chan struct{}
	running          bool
	mu               *sync.Mutex
	logger           *log.Logger
}

func (s *jobsIndexer) Index(input <-chan jobs.Job) (output <-chan jobs.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		panic(errors.New("indexer already running"))
	}
	s.running = true

	realOutput := make(chan jobs.Job)
	go func() {
	loop:
		for {
			select {
			case job, ok := <-input:
				if !ok {
					break loop
				}
				s.logger.Println("Storing job " + job.ID + "...")
				if err := s.storage.UpdateJob(job); err != nil {
					s.logger.Printf("Failed to store job "+job.ID+" (%v).\n", err)
					//TODO retry
				} else {
					realOutput <- job
				}
			case <-s.runContext.Done():
				// TODO log warning that the input did not exit on time.
				break loop
			}
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		close(realOutput)
		close(s.done)
		s.running = false
	}()
	return realOutput
}

func (s *jobsIndexer) Shutdown(ctx context.Context) {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()
	select {
	case <-s.done:
	case <-ctx.Done():
		// Forcefully exit
		s.cancelRunContext()
	}
	<-s.done
}
