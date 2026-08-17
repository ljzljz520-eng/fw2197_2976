package domain

import (
	"errors"
	"sort"
)

type Registry struct {
	Projects map[string]AssessmentProject
}

func NewRegistry() *Registry {
	return &Registry{Projects: make(map[string]AssessmentProject)}
}

func (r *Registry) Add(project AssessmentProject) error {
	if r == nil {
		return errors.New("registry is nil")
	}
	if err := project.Validate(); err != nil {
		return err
	}
	if _, exists := r.Projects[project.ID]; exists {
		return errors.New("project already registered")
	}
	r.Projects[project.ID] = project.Clone()
	return nil
}

func (r *Registry) Replace(project AssessmentProject) error {
	if r == nil {
		return errors.New("registry is nil")
	}
	if err := project.Validate(); err != nil {
		return err
	}
	if _, exists := r.Projects[project.ID]; !exists {
		return ErrNotFound
	}
	r.Projects[project.ID] = project.Clone()
	return nil
}

func (r *Registry) Get(id string) (AssessmentProject, error) {
	if r == nil {
		return AssessmentProject{}, errors.New("registry is nil")
	}
	project, exists := r.Projects[id]
	if !exists {
		return AssessmentProject{}, ErrNotFound
	}
	return project.Clone(), nil
}

func (r *Registry) IDs() []string {
	if r == nil {
		return []string{}
	}
	ids := make([]string, 0, len(r.Projects))
	for id := range r.Projects {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func (r *Registry) Remove(id string) error {
	if r == nil {
		return errors.New("registry is nil")
	}
	if _, exists := r.Projects[id]; !exists {
		return ErrNotFound
	}
	delete(r.Projects, id)
	return nil
}

func (r *Registry) Published() []AssessmentProject {
	projects := make([]AssessmentProject, 0)
	if r == nil {
		return projects
	}
	for _, project := range r.Projects {
		if project.Status == ProjectPublished {
			projects = append(projects, project.Clone())
		}
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].ID < projects[j].ID })
	return projects
}

func (r *Registry) Count(status ProjectStatus) int {
	if r == nil {
		return 0
	}
	count := 0
	for _, project := range r.Projects {
		if status == "" || project.Status == status {
			count++
		}
	}
	return count
}

func (r *Registry) Apply(id string, fn func(*AssessmentProject) error) error {
	if r == nil || fn == nil {
		return errors.New("registry and callback are required")
	}
	project, err := r.Get(id)
	if err != nil {
		return err
	}
	if err := fn(&project); err != nil {
		return err
	}
	return r.Replace(project)
}

func (r *Registry) Transition(id string, status ProjectStatus) error {
	return r.Apply(id, func(project *AssessmentProject) error {
		return project.Transition(status)
	})
}

func (r *Registry) LearnerCount(id string) (int, error) {
	project, err := r.Get(id)
	if err != nil {
		return 0, err
	}
	return len(project.Learners), nil
}

func (r *Registry) Snapshot() map[string]ProjectStatus {
	result := make(map[string]ProjectStatus)
	if r == nil {
		return result
	}
	for id, project := range r.Projects {
		result[id] = project.Status
	}
	return result
}

func (r *Registry) Clear() {
	if r == nil {
		return
	}
	for id := range r.Projects {
		delete(r.Projects, id)
	}
}

func (r *Registry) IsEmpty() bool {
	return r == nil || len(r.Projects) == 0
}

func (r *Registry) Statuses() []ProjectStatus {
	result := make([]ProjectStatus, 0)
	if r == nil {
		return result
	}
	for _, project := range r.Projects {
		result = append(result, project.Status)
	}
	return result
}

func (r *Registry) Size() int {
	if r == nil {
		return 0
	}
	return len(r.Projects)
}
