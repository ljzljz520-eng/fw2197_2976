package workflow

import (
	"errors"
	"sync"

	"skillsassessment/domain"
	"skillsassessment/service"
)

type BatchOutcome struct {
	Added  int
	Errors []error
}

func AddLearnersConcurrently(svc *service.ProjectService, projectID string, learners []domain.Learner) BatchOutcome {
	result := BatchOutcome{Errors: make([]error, 0)}
	if svc == nil {
		result.Errors = append(result.Errors, errors.New("service is required"))
		return result
	}
	var mu sync.Mutex
	var wait sync.WaitGroup
	start := make(chan struct{})
	for _, learner := range learners {
		current := learner
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			err := svc.AddLearner(projectID, current)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Errors = append(result.Errors, err)
				return
			}
			result.Added++
		}()
	}
	close(start)
	wait.Wait()
	return result
}
