package workflow

import (
	"errors"

	"skillsassessment/domain"
	"skillsassessment/service"
)

type ReviewWorkflow struct {
	Policy domain.ReviewPolicy
}

type ReviewWorkflowResult struct {
	Project  domain.AssessmentProject
	Decision domain.PolicyDecision
	Review   service.ReviewResult
}

func NewReviewWorkflow(policy domain.ReviewPolicy) ReviewWorkflow {
	return ReviewWorkflow{Policy: policy}
}

func (w ReviewWorkflow) Execute(svc *service.ProjectService, id, reviewer string) (ReviewWorkflowResult, error) {
	if svc == nil {
		return ReviewWorkflowResult{}, errors.New("service is required")
	}
	project, err := svc.GetProject(id)
	if err != nil {
		return ReviewWorkflowResult{}, err
	}
	decision := w.Policy.Evaluate(project)
	if !decision.Allowed {
		return ReviewWorkflowResult{Project: project, Decision: decision}, errors.New(decision.Reasons[0])
	}
	review, err := svc.ReviewProject(id, reviewer, true, "")
	if err != nil {
		return ReviewWorkflowResult{}, err
	}
	project, err = svc.PublishProject(id, reviewer)
	if err != nil {
		return ReviewWorkflowResult{}, err
	}
	return ReviewWorkflowResult{Project: project, Decision: decision, Review: review}, nil
}

func RejectReview(svc *service.ProjectService, id, reviewer, reason string) error {
	if svc == nil {
		return errors.New("service is required")
	}
	_, err := svc.ReviewProject(id, reviewer, false, reason)
	return err
}
