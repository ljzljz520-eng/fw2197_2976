package service

import (
	"errors"
	"fmt"
	"strings"

	"skillsassessment/domain"
	"skillsassessment/repository"
)

type ProjectService struct {
	projects *repository.ProjectRepository
	learners *repository.LearnerRepository
}

func NewProjectService(projects *repository.ProjectRepository, learners *repository.LearnerRepository) *ProjectService {
	return &ProjectService{projects: projects, learners: learners}
}

func (s *ProjectService) RegisterProject(id, name, creator string, certification domain.Certification) (domain.AssessmentProject, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" || strings.TrimSpace(creator) == "" {
		return domain.AssessmentProject{}, errors.New("project identity is required")
	}
	project := domain.AssessmentProject{ID: id, Name: name, Status: domain.ProjectDraft, Version: 1, Certification: certification}
	project.Audit.CreatedBy = creator
	project.Audit.Events = []string{"REGISTERED"}
	if err := s.projects.Save(project); err != nil {
		return domain.AssessmentProject{}, err
	}
	return project, nil
}

func (s *ProjectService) GetProject(id string) (domain.AssessmentProject, error) {
	return s.projects.Find(id)
}

func (s *ProjectService) PublishProject(id, reviewer string) (domain.AssessmentProject, error) {
	project, err := s.projects.Find(id)
	if err != nil {
		return domain.AssessmentProject{}, err
	}
	if project.Status != domain.ProjectDraft {
		return domain.AssessmentProject{}, domain.ErrInvalidTransition
	}
	if strings.TrimSpace(reviewer) == "" {
		return domain.AssessmentProject{}, errors.New("reviewer is required")
	}
	project.Audit.ReviewedBy = append(project.Audit.ReviewedBy, reviewer)
	if err := project.Transition(domain.ProjectPublished); err != nil {
		return domain.AssessmentProject{}, err
	}
	if err := s.projects.Save(project); err != nil {
		return domain.AssessmentProject{}, err
	}
	return project, nil
}

func (s *ProjectService) ArchiveProject(id, reviewer string) (domain.AssessmentProject, error) {
	project, err := s.projects.Find(id)
	if err != nil {
		return domain.AssessmentProject{}, err
	}
	if project.Status != domain.ProjectPublished {
		return domain.AssessmentProject{}, domain.ErrInvalidTransition
	}
	if strings.TrimSpace(reviewer) == "" {
		return domain.AssessmentProject{}, errors.New("reviewer is required")
	}
	project.Audit.ReviewedBy = append(project.Audit.ReviewedBy, reviewer)
	if err := project.Transition(domain.ProjectArchived); err != nil {
		return domain.AssessmentProject{}, err
	}
	if err := s.projects.Save(project); err != nil {
		return domain.AssessmentProject{}, err
	}
	return project, nil
}

func (s *ProjectService) ReopenProject(id string) (domain.AssessmentProject, error) {
	project, err := s.projects.Find(id)
	if err != nil {
		return domain.AssessmentProject{}, err
	}
	if project.Status == domain.ProjectArchived {
		return domain.AssessmentProject{}, domain.ErrInvalidTransition
	}
	return project, nil
}

func (s *ProjectService) RemoveProject(id string) error {
	project, err := s.projects.Find(id)
	if err != nil {
		return err
	}
	if project.Status == domain.ProjectPublished {
		return fmt.Errorf("cannot remove published project")
	}
	return s.projects.Remove(id)
}

func (s *ProjectService) ListProjects() ([]domain.AssessmentProject, error) {
	return s.projects.All()
}
