package service

import (
	"errors"
	"strings"

	"skillsassessment/domain"
	"skillsassessment/repository"
)

type AuditService struct {
	projects *repository.ProjectRepository
	events   *repository.EventRepository
}

func NewAuditService(projects *repository.ProjectRepository) *AuditService {
	return &AuditService{projects: projects, events: repository.NewEventRepository(projects)}
}

func (s *AuditService) Events(id string) ([]repository.Event, error) {
	if s == nil || s.events == nil {
		return nil, errors.New("audit service is required")
	}
	return s.events.List(id)
}

func (s *AuditService) HasEvent(id, name string) (bool, error) {
	if strings.TrimSpace(name) == "" {
		return false, errors.New("event name is required")
	}
	return s.events.Contains(id, name)
}

func (s *AuditService) CurrentState(id string) (domain.ProjectStatus, error) {
	project, err := s.projects.Find(id)
	if err != nil {
		return "", err
	}
	return project.Status, nil
}

func (s *AuditService) VerifyStateEvent(id string) error {
	project, err := s.projects.Find(id)
	if err != nil {
		return err
	}
	last, err := s.events.Last(id)
	if err != nil {
		return err
	}
	if last.Name != repository.EventForStatus(project.Status) {
		return errors.New("state event mismatch")
	}
	return nil
}

func (s *AuditService) AddManualEvent(id, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("event name is required")
	}
	project, err := s.projects.Find(id)
	if err != nil {
		return err
	}
	if project.Status == domain.ProjectArchived {
		return domain.ErrInvalidTransition
	}
	project.Audit.Events = append(project.Audit.Events, "MANUAL:"+name)
	project.Version++
	return s.projects.Save(project)
}
