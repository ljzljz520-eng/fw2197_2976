package repository

import (
	"fmt"
	"strings"

	"skillsassessment/domain"
)

type EventRepository struct {
	projects *ProjectRepository
}

type Event struct {
	ProjectID string
	Sequence  int
	Name      string
}

func NewEventRepository(projects *ProjectRepository) *EventRepository {
	return &EventRepository{projects: projects}
}

func (r *EventRepository) List(projectID string) ([]Event, error) {
	project, err := r.projects.Find(projectID)
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(project.Audit.Events))
	for index, name := range project.Audit.Events {
		events = append(events, Event{ProjectID: project.ID, Sequence: index + 1, Name: name})
	}
	return events, nil
}

func (r *EventRepository) Contains(projectID, eventName string) (bool, error) {
	events, err := r.List(projectID)
	if err != nil {
		return false, err
	}
	for _, event := range events {
		if event.Name == eventName || strings.HasPrefix(event.Name, eventName+":") {
			return true, nil
		}
	}
	return false, nil
}

func (r *EventRepository) Last(projectID string) (Event, error) {
	events, err := r.List(projectID)
	if err != nil {
		return Event{}, err
	}
	if len(events) == 0 {
		return Event{}, fmt.Errorf("project %s has no events", projectID)
	}
	return events[len(events)-1], nil
}

func EventForStatus(status domain.ProjectStatus) string {
	switch status {
	case domain.ProjectDraft:
		return "REGISTERED"
	case domain.ProjectPublished:
		return "PUBLISHED"
	case domain.ProjectArchived:
		return "ARCHIVED"
	default:
		return "UNKNOWN"
	}
}
