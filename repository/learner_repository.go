package repository

import (
	"sort"

	"skillsassessment/domain"
)

type LearnerRepository struct {
	projects *ProjectRepository
}

func NewLearnerRepository(projects *ProjectRepository) *LearnerRepository {
	return &LearnerRepository{projects: projects}
}

func (r *LearnerRepository) Add(projectID string, learner domain.Learner) error {
	project, err := r.projects.Find(projectID)
	if err != nil {
		return err
	}
	if err := learner.Validate(); err != nil {
		return err
	}
	for _, current := range project.Learners {
		if current.ID == learner.ID {
			return domain.ErrDuplicateLearner
		}
	}
	project.Learners = append(project.Learners, learner)
	project.Version++
	project.Audit.Events = append(project.Audit.Events, "LEARNER_ADDED:"+learner.ID)
	return r.projects.Save(project)
}

func (r *LearnerRepository) Update(projectID string, learner domain.Learner) error {
	project, err := r.projects.Find(projectID)
	if err != nil {
		return err
	}
	if err := learner.Validate(); err != nil {
		return err
	}
	found := false
	for index := range project.Learners {
		if project.Learners[index].ID == learner.ID {
			project.Learners[index] = learner
			found = true
			break
		}
	}
	if !found {
		return domain.ErrNotFoundLearner
	}
	project.Version++
	project.Audit.Events = append(project.Audit.Events, "LEARNER_UPDATED:"+learner.ID)
	return r.projects.Save(project)
}

func (r *LearnerRepository) Remove(projectID, learnerID string) error {
	project, err := r.projects.Find(projectID)
	if err != nil {
		return err
	}
	index := -1
	for i, learner := range project.Learners {
		if learner.ID == learnerID {
			index = i
			break
		}
	}
	if index < 0 {
		return domain.ErrNotFoundLearner
	}
	project.Learners = append(project.Learners[:index], project.Learners[index+1:]...)
	project.Version++
	project.Audit.Events = append(project.Audit.Events, "LEARNER_REMOVED:"+learnerID)
	return r.projects.Save(project)
}

func (r *LearnerRepository) List(projectID string) ([]domain.Learner, error) {
	project, err := r.projects.Find(projectID)
	if err != nil {
		return nil, err
	}
	learners := append([]domain.Learner(nil), project.Learners...)
	sort.Slice(learners, func(i, j int) bool { return learners[i].ID < learners[j].ID })
	return learners, nil
}
