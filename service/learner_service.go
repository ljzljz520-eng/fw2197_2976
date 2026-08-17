package service

import (
	"errors"
	"strings"

	"skillsassessment/domain"
)

func (s *ProjectService) AddLearner(projectID string, learner domain.Learner) error {
	project, err := s.projects.Find(projectID)
	if err != nil {
		return err
	}
	if project.Status == domain.ProjectArchived {
		return domain.ErrInvalidTransition
	}
	return s.learners.Add(projectID, learner)
}

func (s *ProjectService) UpdateLearner(projectID string, learner domain.Learner) error {
	if strings.TrimSpace(learner.ID) == "" {
		return errors.New("learner id is required")
	}
	project, err := s.projects.Find(projectID)
	if err != nil {
		return err
	}
	if project.Status == domain.ProjectArchived {
		return domain.ErrInvalidTransition
	}
	return s.learners.Update(projectID, learner)
}

func (s *ProjectService) RemoveLearner(projectID, learnerID string) error {
	project, err := s.projects.Find(projectID)
	if err != nil {
		return err
	}
	if project.Status == domain.ProjectArchived {
		return domain.ErrInvalidTransition
	}
	return s.learners.Remove(projectID, learnerID)
}

func (s *ProjectService) ListLearners(projectID string) ([]domain.Learner, error) {
	return s.learners.List(projectID)
}
