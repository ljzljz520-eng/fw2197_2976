package repository

import (
	"sort"
	"strings"

	"skillsassessment/domain"
	"skillsassessment/storage"
)

type CatalogRepository struct {
	projects *ProjectRepository
}

type CatalogEntry struct {
	ID     string
	Name   string
	Status domain.ProjectStatus
	Score  float64
}

func NewCatalogRepository(projects *ProjectRepository) *CatalogRepository {
	return &CatalogRepository{projects: projects}
}

func (r *CatalogRepository) Available() ([]CatalogEntry, error) {
	projects, err := r.projects.Search(storage.Query{Status: domain.ProjectPublished})
	if err != nil {
		return nil, err
	}
	entries := make([]CatalogEntry, 0, len(projects))
	for _, project := range projects {
		entries = append(entries, CatalogEntry{ID: project.ID, Name: project.Name, Status: project.Status, Score: domain.SummarizeScores(project.Learners).Average})
	}
	return entries, nil
}

func (r *CatalogRepository) FindByName(name string) ([]CatalogEntry, error) {
	projects, err := r.projects.All()
	if err != nil {
		return nil, err
	}
	entries := make([]CatalogEntry, 0)
	for _, project := range projects {
		if strings.Contains(strings.ToLower(project.Name), strings.ToLower(name)) {
			entries = append(entries, CatalogEntry{ID: project.ID, Name: project.Name, Status: project.Status, Score: domain.SummarizeScores(project.Learners).Average})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries, nil
}

func (r *CatalogRepository) StatusCounts() (map[domain.ProjectStatus]int, error) {
	projects, err := r.projects.All()
	if err != nil {
		return nil, err
	}
	counts := map[domain.ProjectStatus]int{}
	for _, project := range projects {
		counts[project.Status]++
	}
	return counts, nil
}
