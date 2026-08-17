package domain

import "sort"

type ScoreSummary struct {
	Count   int     `json:"count"`
	Highest float64 `json:"highest"`
	Lowest  float64 `json:"lowest"`
	Average float64 `json:"average"`
}

func (l Learner) OverallScore() float64 {
	return (l.Theory + l.Practical + l.Safety) / 3
}

func SummarizeScores(learners []Learner) ScoreSummary {
	if len(learners) == 0 {
		return ScoreSummary{}
	}
	values := make([]float64, 0, len(learners))
	total := 0.0
	for _, learner := range learners {
		score := learner.OverallScore()
		values = append(values, score)
		total += score
	}
	sort.Float64s(values)
	return ScoreSummary{Count: len(values), Lowest: values[0], Highest: values[len(values)-1], Average: total / float64(len(values))}
}

func ClassifyScore(score float64) string {
	switch {
	case score >= 90:
		return "EXCELLENT"
	case score >= 75:
		return "QUALIFIED"
	case score >= 60:
		return "DEVELOPING"
	default:
		return "RETRAIN"
	}
}

func FormatCertification(c Certification) string {
	if c.Kind == CertificationGrade {
		return "Grade: " + c.Grade
	}
	return "Score: " + formatPercent(c.Percent)
}

func formatPercent(value float64) string {
	whole := int(value)
	fraction := int(value*10) % 10
	if fraction == 0 {
		return itoa(whole) + "%"
	}
	return itoa(whole) + "." + itoa(fraction) + "%"
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	buf := make([]byte, 0, 12)
	for value > 0 {
		buf = append([]byte{byte('0' + value%10)}, buf...)
		value /= 10
	}
	if negative {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
