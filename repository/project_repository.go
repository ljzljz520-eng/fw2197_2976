package repository

import (
	"sort"

	"skillsassessment/domain"
	"skillsassessment/storage"
)

type ProjectRepository struct {
	db *storage.Database
}

func NewProjectRepository(db *storage.Database) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Save(project domain.AssessmentProject) error {
	return r.db.PutProject(project)
}

func (r *ProjectRepository) Find(id string) (domain.AssessmentProject, error) {
	return r.db.GetProject(id)
}

func (r *ProjectRepository) Remove(id string) error {
	return r.db.DeleteProject(id)
}

func (r *ProjectRepository) All() ([]domain.AssessmentProject, error) {
	projects, err := r.db.ListProjects()
	if err != nil {
		return nil, err
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].ID < projects[j].ID })
	return projects, nil
}

func (r *ProjectRepository) Search(query storage.Query) ([]domain.AssessmentProject, error) {
	projects, err := r.db.QueryProjects(query)
	if err != nil {
		return nil, err
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].Name < projects[j].Name })
	return projects, nil
}
