package workflow

import (
	"errors"
	"testing"

	"skillsassessment/domain"
)

func TestWorkflowPublication(t *testing.T) {
	svc := workflowService(t)
	if _, err := RegisterAssessment(svc, "WF-PUB", "Publication", "trainer", domain.Certification{Kind: domain.CertificationScore, Percent: 91}, nil); err != nil {
		t.Fatal(err)
	}
	result, err := PublishAssessment(svc, "WF-PUB", "auditor")
	if err != nil || result.Project.Status != domain.ProjectPublished {
		t.Fatalf("publish workflow: %v %+v", err, result)
	}
	if _, err := ArchiveAssessment(svc, "WF-PUB", "auditor"); err != nil {
		t.Fatal(err)
	}
	project, err := svc.GetProject("WF-PUB")
	if err != nil || project.Status != domain.ProjectArchived {
		t.Fatalf("archive workflow: %v %+v", err, project)
	}
	_, err = PublishArchivedAssessment(svc, "WF-PUB", "auditor")
	if !errors.Is(err, domain.ErrInvalidTransition) {
		t.Fatalf("archived project must reject publish, got %v", err)
	}
}
