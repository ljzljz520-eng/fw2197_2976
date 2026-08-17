package domain

import "math"

type ScoreBand struct {
	Name  string  `json:"name"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int     `json:"count"`
}

type DimensionSummary struct {
	Theory    ScoreSummary `json:"theory"`
	Practical ScoreSummary `json:"practical"`
	Safety    ScoreSummary `json:"safety"`
}

func SummarizeDimensions(learners []Learner) DimensionSummary {
	theory := make([]Learner, 0, len(learners))
	practical := make([]Learner, 0, len(learners))
	safety := make([]Learner, 0, len(learners))
	for _, learner := range learners {
		theory = append(theory, Learner{Theory: learner.Theory, Practical: learner.Theory, Safety: learner.Theory})
		practical = append(practical, Learner{Theory: learner.Practical, Practical: learner.Practical, Safety: learner.Practical})
		safety = append(safety, Learner{Theory: learner.Safety, Practical: learner.Safety, Safety: learner.Safety})
	}
	return DimensionSummary{Theory: summarizeDimension(theory, 0), Practical: summarizeDimension(practical, 1), Safety: summarizeDimension(safety, 2)}
}

func summarizeDimension(learners []Learner, dimension int) ScoreSummary {
	if len(learners) == 0 {
		return ScoreSummary{}
	}
	values := make([]float64, 0, len(learners))
	total := 0.0
	for _, learner := range learners {
		value := learner.Theory
		if dimension == 1 {
			value = learner.Practical
		} else if dimension == 2 {
			value = learner.Safety
		}
		values = append(values, value)
		total += value
	}
	lowest, highest := values[0], values[0]
	for _, value := range values[1:] {
		if value < lowest {
			lowest = value
		}
		if value > highest {
			highest = value
		}
	}
	return ScoreSummary{Count: len(values), Lowest: lowest, Highest: highest, Average: total / float64(len(values))}
}

func BuildBands(learners []Learner) []ScoreBand {
	bands := []ScoreBand{{Name: "EXCELLENT", Min: 90, Max: 100}, {Name: "QUALIFIED", Min: 75, Max: 89.99}, {Name: "DEVELOPING", Min: 60, Max: 74.99}, {Name: "RETRAIN", Min: 0, Max: 59.99}}
	for _, learner := range learners {
		score := learner.OverallScore()
		for index := range bands {
			if score >= bands[index].Min && score <= bands[index].Max {
				bands[index].Count++
				break
			}
		}
	}
	return bands
}

func RoundScore(value float64) float64 {
	return math.Round(value*100) / 100
}

func IsPassing(learner Learner, threshold float64) bool {
	if threshold <= 0 {
		threshold = 60
	}
	return learner.OverallScore() >= threshold && learner.Safety >= threshold
}
