package reporting

import (
	"testing"

	"skillsassessment/domain"
)

func TestReportAndRanking(t *testing.T) {
	project := domain.AssessmentProject{ID: "REP-1", Name: "Report", Status: domain.ProjectPublished, Certification: domain.Certification{Kind: domain.CertificationScore, Percent: 88}, Learners: []domain.Learner{{ID: "1", Name: "Top", Gender: "F", Email: "top@example.test", Theory: 95, Practical: 95, Safety: 95}, {ID: "2", Name: "Next", Gender: "M", Email: "next@example.test", Theory: 70, Practical: 70, Safety: 70}}}
	report := BuildProjectReport(project)
	if report.Summary.Highest != 95 || report.Breakdown["EXCELLENT"] != 1 || report.Certification != "Score: 88%" {
		t.Fatalf("unexpected report: %+v", report)
	}
	ranked := TopLearners(project.Learners, 1)
	if len(ranked) != 1 || ranked[0].Learner.ID != "1" || ranked[0].Rank != 1 {
		t.Fatalf("unexpected ranking: %+v", ranked)
	}
}
