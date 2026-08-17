package domain

import "testing"

func TestScoreSummaryAndClassification(t *testing.T) {
	learners := []Learner{{ID: "A", Name: "A", Gender: "M", Email: "a@example.test", Theory: 90, Practical: 90, Safety: 90}, {ID: "B", Name: "B", Gender: "F", Email: "b@example.test", Theory: 60, Practical: 60, Safety: 60}}
	summary := SummarizeScores(learners)
	if summary.Count != 2 || summary.Highest != 90 || summary.Lowest != 60 || summary.Average != 75 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if ClassifyScore(95) != "EXCELLENT" || ClassifyScore(75) != "QUALIFIED" || ClassifyScore(60) != "DEVELOPING" || ClassifyScore(30) != "RETRAIN" {
		t.Fatal("classification mismatch")
	}
	if FormatCertification(Certification{Kind: CertificationScore, Percent: 88.5}) != "Score: 88.5%" {
		t.Fatal("score formatting mismatch")
	}
}

func TestFilterLearners(t *testing.T) {
	learners := []Learner{{ID: "A", Name: "Alice", Gender: "F", Email: "a@example.test", Theory: 90, Practical: 90, Safety: 90}, {ID: "B", Name: "Bob", Gender: "M", Email: "b@example.test", Theory: 50, Practical: 50, Safety: 50}}
	result := FilterLearners(learners, LearnerFilter{Gender: "F", MinOverall: 80})
	if len(result) != 1 || result[0].ID != "A" {
		t.Fatalf("unexpected filter result: %+v", result)
	}
}
