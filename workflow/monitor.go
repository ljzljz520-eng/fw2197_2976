package workflow

import (
	"errors"
	"sort"

	"skillsassessment/domain"
	"skillsassessment/service"
)

type Monitor struct {
	Service *service.ProjectService
}

type MonitorRow struct {
	ID       string
	Status   domain.ProjectStatus
	Learners int
	Average  float64
}

func (m Monitor) Snapshot() ([]MonitorRow, error) {
	if m.Service == nil {
		return nil, errors.New("service is required")
	}
	projects, err := m.Service.ListProjects()
	if err != nil {
		return nil, err
	}
	rows := make([]MonitorRow, 0, len(projects))
	for _, project := range projects {
		rows = append(rows, MonitorRow{ID: project.ID, Status: project.Status, Learners: len(project.Learners), Average: domain.SummarizeScores(project.Learners).Average})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
	return rows, nil
}

func (m Monitor) CountByStatus() (map[domain.ProjectStatus]int, error) {
	rows, err := m.Snapshot()
	if err != nil {
		return nil, err
	}
	counts := map[domain.ProjectStatus]int{}
	for _, row := range rows {
		counts[row.Status]++
	}
	return counts, nil
}

func (m Monitor) NeedsAttention(row MonitorRow, threshold float64) bool {
	if row.Status == domain.ProjectArchived {
		return false
	}
	if row.Learners == 0 {
		return true
	}
	return row.Average < threshold
}
