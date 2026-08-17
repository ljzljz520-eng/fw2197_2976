package reporting

import (
	"sort"

	"skillsassessment/domain"
)

type Insight struct {
	Code     string
	Severity string
	Message  string
}

func GenerateInsights(project domain.AssessmentProject) []Insight {
	insights := make([]Insight, 0)
	if len(project.Learners) == 0 {
		return append(insights, Insight{Code: "NO_LEARNERS", Severity: "HIGH", Message: "project has no learners"})
	}
	summary := domain.SummarizeScores(project.Learners)
	if summary.Average < 60 {
		insights = append(insights, Insight{Code: "LOW_AVERAGE", Severity: "HIGH", Message: "overall average needs intervention"})
	} else if summary.Average < 75 {
		insights = append(insights, Insight{Code: "WATCH_AVERAGE", Severity: "MEDIUM", Message: "overall average needs monitoring"})
	}
	dimensions := domain.SummarizeDimensions(project.Learners)
	if dimensions.Safety.Average < dimensions.Theory.Average {
		insights = append(insights, Insight{Code: "SAFETY_GAP", Severity: "HIGH", Message: "safety trails theory"})
	}
	if project.Status == domain.ProjectArchived {
		insights = append(insights, Insight{Code: "ARCHIVED", Severity: "INFO", Message: "project is archived"})
	}
	sort.Slice(insights, func(i, j int) bool { return insights[i].Code < insights[j].Code })
	return insights
}

func CountSeverity(insights []Insight, severity string) int {
	count := 0
	for _, insight := range insights {
		if insight.Severity == severity {
			count++
		}
	}
	return count
}

func CriticalInsight(insights []Insight) (Insight, bool) {
	for _, insight := range insights {
		if insight.Severity == "HIGH" {
			return insight, true
		}
	}
	return Insight{}, false
}
