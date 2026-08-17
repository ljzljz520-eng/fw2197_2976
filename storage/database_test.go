package storage

import (
	"path/filepath"
	"testing"

	"skillsassessment/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "skills.db")
	project := domain.AssessmentProject{ID: "PERSIST-1", Name: "Reopen", Status: domain.ProjectDraft, Version: 1, Certification: domain.Certification{Kind: domain.CertificationScore, Percent: 82}, Learners: []domain.Learner{{ID: "L-1", Name: "A", Gender: "F", Email: "a@example.test", Theory: 80, Practical: 82, Safety: 84}}}
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.PutProject(project); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	loaded, err := db.GetProject(project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != project.ID || loaded.Certification.Percent != 82 || len(loaded.Learners) != 1 {
		t.Fatalf("persistence mismatch: %+v", loaded)
	}
}

func TestDatabaseQueryAndSnapshot(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, id := range []string{"A-1", "B-1"} {
		project := domain.AssessmentProject{ID: id, Name: id, Status: domain.ProjectDraft, Version: 1, Certification: domain.Certification{Kind: domain.CertificationGrade, Grade: "PASS"}}
		if err := db.PutProject(project); err != nil {
			t.Fatal(err)
		}
	}
	projects, err := db.QueryProjects(Query{Prefix: "A-"})
	if err != nil || len(projects) != 1 {
		t.Fatalf("query mismatch: %v %+v", err, projects)
	}
	snapshot, err := db.Snapshot()
	if err != nil || snapshot.Projects != 2 || snapshot.Audits != 2 {
		t.Fatalf("snapshot mismatch: %v %+v", err, snapshot)
	}
}
