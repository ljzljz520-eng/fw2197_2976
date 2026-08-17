package domain

import "testing"

func TestLearnerValidationAndTransition(t *testing.T) {
	learner := Learner{ID: "L1", Name: "Lin", Gender: "F", Email: "lin@example.test", Theory: 80, Practical: 90, Safety: 70}
	if err := learner.Validate(); err != nil {
		t.Fatal(err)
	}
	project := AssessmentProject{ID: "P1", Name: "Electrical", Status: ProjectDraft, Version: 1, Certification: Certification{Kind: CertificationGrade, Grade: "PASS"}}
	if err := project.Transition(ProjectPublished); err != nil {
		t.Fatal(err)
	}
	if project.Status != ProjectPublished || project.Version != 2 {
		t.Fatalf("unexpected transition: %+v", project)
	}
	if err := project.Transition(ProjectPublished); err == nil {
		t.Fatal("expected invalid transition")
	}
}

func TestCertificationValidation(t *testing.T) {
	cases := []Certification{{Kind: CertificationGrade, Grade: "PASS"}, {Kind: CertificationScore, Percent: 88}}
	for _, current := range cases {
		if err := current.Validate(); err != nil {
			t.Fatalf("%+v: %v", current, err)
		}
	}
	if err := (Certification{Kind: CertificationGrade, Grade: "MAYBE"}).Validate(); err == nil {
		t.Fatal("expected invalid grade")
	}
}
