package workflow

import (
	"path/filepath"
	"testing"

	"skillsassessment/domain"
	"skillsassessment/repository"
	"skillsassessment/service"
	"skillsassessment/storage"
)

func workflowService(t *testing.T) *service.ProjectService {
	t.Helper()
	db, err := storage.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	projects := repository.NewProjectRepository(db)
	return service.NewProjectService(projects, repository.NewLearnerRepository(projects))
}

func TestWorkflowRegistration(t *testing.T) {
	svc := workflowService(t)
	result, err := RegisterAssessment(svc, "WF-REG", "Registration", "trainer", domain.Certification{Kind: domain.CertificationGrade, Grade: "PASS"}, []domain.Learner{{ID: "L1", Name: "One", Gender: "F", Email: "one@example.test", Theory: 80, Practical: 85, Safety: 90}})
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateRegistration(result); err != nil || result.Learners != 1 {
		t.Fatalf("invalid registration: %v %+v", err, result)
	}
}
