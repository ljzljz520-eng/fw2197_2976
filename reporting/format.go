package reporting

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"skillsassessment/domain"
)

func FormatReport(report ProjectReport) string {
	parts := []string{report.ProjectID, report.ProjectName, string(report.Status), fmt.Sprintf("average=%.2f", report.Summary.Average), report.Certification}
	keys := make([]string, 0, len(report.Breakdown))
	for key := range report.Breakdown {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, key+"="+itoa(report.Breakdown[key]))
	}
	return strings.Join(parts, " | ")
}

func MarshalReport(report ProjectReport) ([]byte, error) {
	return json.Marshal(report)
}

func MarshalReports(reports []ProjectReport) ([]byte, error) {
	return json.Marshal(SortReports(reports))
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	result := ""
	for value > 0 {
		result = string(rune('0'+value%10)) + result
		value /= 10
	}
	return result
}

func StatusLabel(status domain.ProjectStatus) string {
	switch status {
	case domain.ProjectDraft:
		return "Draft"
	case domain.ProjectPublished:
		return "Available"
	case domain.ProjectArchived:
		return "Archived"
	default:
		return "Unknown"
	}
}
