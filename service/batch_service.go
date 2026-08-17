package service

import (
	"errors"
	"strings"

	"skillsassessment/domain"
)

type BatchService struct {
	project *ProjectService
}

type BatchRequest struct {
	ProjectID string
	Learners  []domain.Learner
}

type BatchResult struct {
	Accepted []string
	Rejected map[string]string
}

func NewBatchService(project *ProjectService) *BatchService {
	return &BatchService{project: project}
}

func (s *BatchService) Apply(request BatchRequest) (BatchResult, error) {
	if s == nil || s.project == nil {
		return BatchResult{}, errors.New("project service is required")
	}
	if strings.TrimSpace(request.ProjectID) == "" {
		return BatchResult{}, errors.New("project id is required")
	}
	result := BatchResult{Accepted: make([]string, 0), Rejected: map[string]string{}}
	for _, learner := range request.Learners {
		if err := s.project.AddLearner(request.ProjectID, learner); err != nil {
			result.Rejected[learner.ID] = err.Error()
			continue
		}
		result.Accepted = append(result.Accepted, learner.ID)
	}
	return result, nil
}

func (s *BatchService) Validate(request BatchRequest) []error {
	errorsFound := make([]error, 0)
	seen := map[string]bool{}
	for _, learner := range request.Learners {
		if seen[learner.ID] {
			errorsFound = append(errorsFound, domain.ErrDuplicateLearner)
			continue
		}
		seen[learner.ID] = true
		if err := learner.Validate(); err != nil {
			errorsFound = append(errorsFound, err)
		}
	}
	return errorsFound
}

func (s *BatchService) AcceptedCount(result BatchResult) int {
	return len(result.Accepted)
}
