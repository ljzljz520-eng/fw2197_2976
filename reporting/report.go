package reporting

import (
	"sort"

	"skillsassessment/domain"
)

type ProjectReport struct {
	ProjectID     string
	ProjectName   string
	Status        domain.ProjectStatus
	Summary       domain.ScoreSummary
	Breakdown     map[string]int
	Certification string
}

func BuildProjectReport(project domain.AssessmentProject) ProjectReport {
	breakdown := map[string]int{}
	for _, learner := range project.Learners {
		class := domain.ClassifyScore(learner.OverallScore())
		breakdown[class]++
	}
	return ProjectReport{ProjectID: project.ID, ProjectName: project.Name, Status: project.Status, Summary: domain.SummarizeScores(project.Learners), Breakdown: breakdown, Certification: domain.FormatCertification(project.Certification)}
}

func SortReports(reports []ProjectReport) []ProjectReport {
	copy := append([]ProjectReport(nil), reports...)
	sort.Slice(copy, func(i, j int) bool {
		if copy[i].Summary.Average == copy[j].Summary.Average {
			return copy[i].ProjectID < copy[j].ProjectID
		}
		return copy[i].Summary.Average > copy[j].Summary.Average
	})
	return copy
}
