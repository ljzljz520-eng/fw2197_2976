package workflow

import (
	"errors"

	"skillsassessment/domain"
	"skillsassessment/service"
)

type PublicationResult struct {
	Project domain.AssessmentProject
	Review  service.ReviewResult
}

func PublishAssessment(svc *service.ProjectService, id, reviewer string) (PublicationResult, error) {
	if svc == nil {
		return PublicationResult{}, errors.New("service is required")
	}
	review, err := svc.ReviewProject(id, reviewer, true, "")
	if err != nil {
		return PublicationResult{}, err
	}
	project, err := svc.PublishProject(id, reviewer)
	if err != nil {
		return PublicationResult{}, err
	}
	return PublicationResult{Project: project, Review: review}, nil
}

func ArchiveAssessment(svc *service.ProjectService, id, reviewer string) (domain.AssessmentProject, error) {
	if svc == nil {
		return domain.AssessmentProject{}, errors.New("service is required")
	}
	return svc.ArchiveProject(id, reviewer)
}

func PublishArchivedAssessment(svc *service.ProjectService, id, reviewer string) (domain.AssessmentProject, error) {
	if svc == nil {
		return domain.AssessmentProject{}, errors.New("service is required")
	}
	return svc.PublishProject(id, reviewer)
}
