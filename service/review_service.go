package service

import (
	"errors"
	"strings"

	"skillsassessment/domain"
)

type ReviewResult struct {
	ProjectID string
	Reviewer  string
	Approved  bool
	Reason    string
}

func (s *ProjectService) ReviewProject(id, reviewer string, approve bool, reason string) (ReviewResult, error) {
	if strings.TrimSpace(reviewer) == "" {
		return ReviewResult{}, errors.New("reviewer is required")
	}
	project, err := s.projects.Find(id)
	if err != nil {
		return ReviewResult{}, err
	}
	if project.Status == domain.ProjectArchived {
		return ReviewResult{}, domain.ErrInvalidTransition
	}
	if !approve && strings.TrimSpace(reason) == "" {
		return ReviewResult{}, errors.New("rejection reason is required")
	}
	project.Audit.ReviewedBy = append(project.Audit.ReviewedBy, reviewer)
	project.Audit.Events = append(project.Audit.Events, "REVIEWED:"+reviewer)
	if err := s.projects.Save(project); err != nil {
		return ReviewResult{}, err
	}
	return ReviewResult{ProjectID: id, Reviewer: reviewer, Approved: approve, Reason: reason}, nil
}

func (s *ProjectService) AuditEvents(id string) ([]string, error) {
	project, err := s.projects.Find(id)
	if err != nil {
		return nil, err
	}
	return append([]string(nil), project.Audit.Events...), nil
}
