package domain

import "strings"

type LearnerFilter struct {
	Name       string
	Gender     string
	MinOverall float64
	MaxOverall float64
	Class      string
}

func (f LearnerFilter) Match(l Learner) bool {
	if f.Name != "" && !strings.Contains(strings.ToLower(l.Name), strings.ToLower(f.Name)) {
		return false
	}
	if f.Gender != "" && !strings.EqualFold(l.Gender, f.Gender) {
		return false
	}
	score := l.OverallScore()
	if f.MinOverall > 0 && score < f.MinOverall {
		return false
	}
	if f.MaxOverall > 0 && score > f.MaxOverall {
		return false
	}
	if f.Class != "" && ClassifyScore(score) != f.Class {
		return false
	}
	return true
}

func FilterLearners(learners []Learner, filter LearnerFilter) []Learner {
	matched := make([]Learner, 0)
	for _, learner := range learners {
		if filter.Match(learner) {
			matched = append(matched, learner)
		}
	}
	return matched
}
