package workflow

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"

	"skillsassessment/domain"
	"skillsassessment/service"
)

func ExportAssessment(svc *service.ProjectService, id string, writer io.Writer) error {
	if svc == nil || writer == nil {
		return errors.New("service and writer are required")
	}
	project, err := svc.GetProject(id)
	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"id", "name", "gender", "email", "theory", "practical", "safety", "overall", "classification", "certification"}); err != nil {
		return err
	}
	for _, learner := range project.Learners {
		if err := csvWriter.Write([]string{learner.ID, learner.Name, learner.Gender, learner.Email, formatScore(learner.Theory), formatScore(learner.Practical), formatScore(learner.Safety), formatScore(learner.OverallScore()), domain.ClassifyScore(learner.OverallScore()), domain.FormatCertification(project.Certification)}); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

func formatScore(score float64) string {
	return fmt.Sprintf("%.2f", score)
}

func ExportSummary(svc *service.ProjectService, id string, writer io.Writer) error {
	if svc == nil || writer == nil {
		return errors.New("service and writer are required")
	}
	project, err := svc.GetProject(id)
	if err != nil {
		return err
	}
	summary := domain.SummarizeScores(project.Learners)
	_, err = io.WriteString(writer, "count="+strconv.Itoa(summary.Count)+" highest="+formatScore(summary.Highest)+" lowest="+formatScore(summary.Lowest)+" average="+formatScore(summary.Average)+"\n")
	return err
}
