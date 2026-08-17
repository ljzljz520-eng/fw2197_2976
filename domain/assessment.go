package domain

import (
	"errors"
	"fmt"
	"strings"
)

type ProjectStatus string

const (
	ProjectDraft     ProjectStatus = "DRAFT"
	ProjectPublished ProjectStatus = "PUBLISHED"
	ProjectArchived  ProjectStatus = "ARCHIVED"
)

type CertificationKind string

const (
	CertificationGrade CertificationKind = "GRADE"
	CertificationScore CertificationKind = "SCORE"
)

type Certification struct {
	Kind    CertificationKind `json:"kind"`
	Grade   string            `json:"grade,omitempty"`
	Percent float64           `json:"percent,omitempty"`
}

type Learner struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Gender    string  `json:"gender"`
	Email     string  `json:"email"`
	Theory    float64 `json:"theory"`
	Practical float64 `json:"practical"`
	Safety    float64 `json:"safety"`
}

type AssessmentProject struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Status        ProjectStatus `json:"status"`
	Learners      []Learner     `json:"learners"`
	Certification Certification `json:"certification"`
	Audit         AuditTrail    `json:"audit"`
	Version       int           `json:"version"`
}

type AuditTrail struct {
	CreatedBy  string   `json:"created_by"`
	ReviewedBy []string `json:"reviewed_by"`
	Events     []string `json:"events"`
}

var (
	ErrInvalidProject    = errors.New("invalid assessment project")
	ErrInvalidLearner    = errors.New("invalid learner")
	ErrInvalidTransition = errors.New("invalid project state transition")
	ErrDuplicateLearner  = errors.New("duplicate learner")
)

func (l Learner) Validate() error {
	if strings.TrimSpace(l.ID) == "" || strings.TrimSpace(l.Name) == "" {
		return ErrInvalidLearner
	}
	if !strings.Contains(l.Email, "@") || strings.TrimSpace(l.Gender) == "" {
		return ErrInvalidLearner
	}
	for _, score := range []float64{l.Theory, l.Practical, l.Safety} {
		if score < 0 || score > 100 {
			return ErrInvalidLearner
		}
	}
	return nil
}

func (p AssessmentProject) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Name) == "" {
		return ErrInvalidProject
	}
	if p.Status != ProjectDraft && p.Status != ProjectPublished && p.Status != ProjectArchived {
		return ErrInvalidProject
	}
	if p.Version < 1 {
		return ErrInvalidProject
	}
	seen := make(map[string]bool)
	for _, learner := range p.Learners {
		if err := learner.Validate(); err != nil {
			return err
		}
		if seen[learner.ID] {
			return ErrDuplicateLearner
		}
		seen[learner.ID] = true
	}
	if err := p.Certification.Validate(); err != nil {
		return err
	}
	return nil
}

func (c Certification) Validate() error {
	switch c.Kind {
	case CertificationGrade:
		if c.Grade != "PASS" && c.Grade != "FAIL" {
			return fmt.Errorf("%w: grade", ErrInvalidProject)
		}
	case CertificationScore:
		if c.Percent < 0 || c.Percent > 100 {
			return fmt.Errorf("%w: percent", ErrInvalidProject)
		}
	default:
		return fmt.Errorf("%w: certification kind", ErrInvalidProject)
	}
	return nil
}

func (p AssessmentProject) CanTransition(to ProjectStatus) bool {
	switch p.Status {
	case ProjectDraft:
		return to == ProjectPublished
	case ProjectPublished:
		return to == ProjectArchived
	case ProjectArchived:
		return false
	default:
		return false
	}
}

func (p *AssessmentProject) Transition(to ProjectStatus) error {
	if !p.CanTransition(to) {
		return ErrInvalidTransition
	}
	p.Status = to
	p.Version++
	p.Audit.Events = append(p.Audit.Events, string(to))
	return nil
}

func (p AssessmentProject) Clone() AssessmentProject {
	copy := p
	copy.Learners = append([]Learner(nil), p.Learners...)
	copy.Audit.ReviewedBy = append([]string(nil), p.Audit.ReviewedBy...)
	copy.Audit.Events = append([]string(nil), p.Audit.Events...)
	return copy
}
