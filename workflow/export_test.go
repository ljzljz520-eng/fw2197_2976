package workflow

import (
	"bytes"
	"strings"
	"testing"

	"skillsassessment/domain"
)

func TestWorkflowExport(t *testing.T) {
	svc := workflowService(t)
	if _, err := RegisterAssessment(svc, "WF-EXP", "Export", "trainer", domain.Certification{Kind: domain.CertificationGrade, Grade: "FAIL"}, []domain.Learner{{ID: "L1", Name: "Exported", Gender: "M", Email: "e@example.test", Theory: 60, Practical: 65, Safety: 70}}); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := ExportAssessment(svc, "WF-EXP", &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Exported") || !strings.Contains(output.String(), "Grade: FAIL") {
		t.Fatalf("unexpected export: %s", output.String())
	}
	output.Reset()
	if err := ExportSummary(svc, "WF-EXP", &output); err != nil || !strings.Contains(output.String(), "average=") {
		t.Fatalf("unexpected summary: %v %s", err, output.String())
	}
}
