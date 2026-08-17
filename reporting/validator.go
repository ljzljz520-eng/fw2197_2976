package reporting

import (
	"errors"
	"strings"

	"skillsassessment/domain"
)

type ValidationIssue struct {
	Field   string
	Message string
}

func ValidateForExport(project domain.AssessmentProject) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	if strings.TrimSpace(project.Name) == "" {
		issues = append(issues, ValidationIssue{Field: "name", Message: "name is required"})
	}
	if len(project.Learners) == 0 {
		issues = append(issues, ValidationIssue{Field: "learners", Message: "at least one learner is required"})
	}
	for index, learner := range project.Learners {
		if err := learner.Validate(); err != nil {
			issues = append(issues, ValidationIssue{Field: "learners[" + strconvIndex(index) + "]", Message: err.Error()})
		}
	}
	return issues
}

func EnsureExportable(project domain.AssessmentProject) error {
	issues := ValidateForExport(project)
	if len(issues) == 0 {
		return nil
	}
	return errors.New(issues[0].Field + ": " + issues[0].Message)
}

func strconvIndex(index int) string {
	if index == 0 {
		return "0"
	}
	result := ""
	for index > 0 {
		result = string(rune('0'+index%10)) + result
		index /= 10
	}
	return result
}
