package service

import (
	"path/filepath"
	"testing"

	"skillsassessment/domain"
	"skillsassessment/repository"
	"skillsassessment/storage"
)

func newService(t *testing.T) *ProjectService {
	t.Helper()
	db, err := storage.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	projects := repository.NewProjectRepository(db)
	return NewProjectService(projects, repository.NewLearnerRepository(projects))
}

func TestServiceLifecycle(t *testing.T) {
	svc := newService(t)
	project, err := svc.RegisterProject("S1", "Service", "owner", domain.Certification{Kind: domain.CertificationGrade, Grade: "PASS"})
	if err != nil || project.Status != domain.ProjectDraft {
		t.Fatalf("register: %v %+v", err, project)
	}
	if err := svc.AddLearner("S1", domain.Learner{ID: "L1", Name: "Learner", Gender: "F", Email: "l@example.test", Theory: 80, Practical: 80, Safety: 80}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ReviewProject("S1", "reviewer", false, "missing docs"); err != nil {
		t.Fatal(err)
	}
	project, err = svc.PublishProject("S1", "reviewer")
	if err != nil || project.Status != domain.ProjectPublished {
		t.Fatalf("publish: %v %+v", err, project)
	}
	project, err = svc.ArchiveProject("S1", "reviewer")
	if err != nil || project.Status != domain.ProjectArchived {
		t.Fatalf("archive: %v %+v", err, project)
	}
}
