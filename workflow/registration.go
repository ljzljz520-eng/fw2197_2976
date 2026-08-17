package workflow

import (
	"errors"
	"strings"

	"skillsassessment/domain"
	"skillsassessment/service"
)

type RegistrationResult struct {
	Project  domain.AssessmentProject
	Learners int
}

func RegisterAssessment(svc *service.ProjectService, id, name, creator string, certification domain.Certification, learners []domain.Learner) (RegistrationResult, error) {
	if svc == nil {
		return RegistrationResult{}, errors.New("service is required")
	}
	if strings.TrimSpace(id) == "" {
		return RegistrationResult{}, errors.New("project id is required")
	}
	project, err := svc.RegisterProject(id, name, creator, certification)
	if err != nil {
		return RegistrationResult{}, err
	}
	for _, learner := range learners {
		if err := svc.AddLearner(id, learner); err != nil {
			return RegistrationResult{}, err
		}
	}
	project, err = svc.GetProject(id)
	if err != nil {
		return RegistrationResult{}, err
	}
	return RegistrationResult{Project: project, Learners: len(project.Learners)}, nil
}

func ValidateRegistration(result RegistrationResult) error {
	if result.Project.ID == "" || result.Learners != len(result.Project.Learners) {
		return errors.New("registration result is inconsistent")
	}
	return result.Project.Validate()
}
