package reporting

import (
	"sort"

	"skillsassessment/domain"
)

type LearnerRank struct {
	Learner domain.Learner
	Rank    int
	Class   string
}

func RankLearners(learners []domain.Learner) []LearnerRank {
	ordered := append([]domain.Learner(nil), learners...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].OverallScore() > ordered[j].OverallScore() })
	ranked := make([]LearnerRank, 0, len(ordered))
	lastScore := -1.0
	rank := 0
	for index, learner := range ordered {
		if learner.OverallScore() != lastScore {
			rank = index + 1
			lastScore = learner.OverallScore()
		}
		ranked = append(ranked, LearnerRank{Learner: learner, Rank: rank, Class: domain.ClassifyScore(learner.OverallScore())})
	}
	return ranked
}

func TopLearners(learners []domain.Learner, limit int) []LearnerRank {
	if limit < 1 {
		return []LearnerRank{}
	}
	ranked := RankLearners(learners)
	if limit >= len(ranked) {
		return ranked
	}
	return ranked[:limit]
}
