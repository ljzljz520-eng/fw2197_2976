package repository

import (
	"path/filepath"
	"testing"

	"skillsassessment/domain"
	"skillsassessment/storage"
)

func TestProjectAndLearnerRepositories(t *testing.T) {
	db, err := storage.Open(filepath.Join(t.TempDir(), "repo.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	projects := NewProjectRepository(db)
	learners := NewLearnerRepository(projects)
	project := domain.AssessmentProject{ID: "R1", Name: "Repo", Status: domain.ProjectDraft, Version: 1, Certification: domain.Certification{Kind: domain.CertificationGrade, Grade: "PASS"}}
	if err := projects.Save(project); err != nil {
		t.Fatal(err)
	}
	learner := domain.Learner{ID: "L1", Name: "Repo Learner", Gender: "M", Email: "repo@example.test", Theory: 70, Practical: 71, Safety: 72}
	if err := learners.Add("R1", learner); err != nil {
		t.Fatal(err)
	}
	if err := learners.Add("R1", learner); err != domain.ErrDuplicateLearner {
		t.Fatalf("expected duplicate error, got %v", err)
	}
	loaded, err := learners.List("R1")
	if err != nil || len(loaded) != 1 {
		t.Fatalf("unexpected learners: %v %+v", err, loaded)
	}
}
