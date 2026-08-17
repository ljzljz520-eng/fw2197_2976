package domain

import "sort"

type TrendPoint struct {
	Label   string  `json:"label"`
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}

type Analytics struct {
	ProjectID string       `json:"project_id"`
	Bands     []ScoreBand  `json:"bands"`
	Trend     []TrendPoint `json:"trend"`
	PassRate  float64      `json:"pass_rate"`
}

func AnalyzeProject(project AssessmentProject) Analytics {
	bands := BuildBands(project.Learners)
	trend := make([]TrendPoint, 0, 3)
	if len(project.Learners) > 0 {
		dimensions := SummarizeDimensions(project.Learners)
		trend = append(trend, TrendPoint{Label: "theory", Average: RoundScore(dimensions.Theory.Average), Count: dimensions.Theory.Count})
		trend = append(trend, TrendPoint{Label: "practical", Average: RoundScore(dimensions.Practical.Average), Count: dimensions.Practical.Count})
		trend = append(trend, TrendPoint{Label: "safety", Average: RoundScore(dimensions.Safety.Average), Count: dimensions.Safety.Count})
	}
	passing := 0
	for _, learner := range project.Learners {
		if IsPassing(learner, 60) {
			passing++
		}
	}
	passRate := 0.0
	if len(project.Learners) > 0 {
		passRate = RoundScore(float64(passing) / float64(len(project.Learners)) * 100)
	}
	return Analytics{ProjectID: project.ID, Bands: bands, Trend: trend, PassRate: passRate}
}

func MergeAnalytics(items []Analytics) Analytics {
	merged := Analytics{Bands: BuildBands(nil), Trend: make([]TrendPoint, 0)}
	if len(items) == 0 {
		return merged
	}
	counts := make(map[string]int)
	totals := make(map[string]float64)
	for _, item := range items {
		if merged.ProjectID == "" {
			merged.ProjectID = item.ProjectID
		}
		for _, band := range item.Bands {
			for index := range merged.Bands {
				if merged.Bands[index].Name == band.Name {
					merged.Bands[index].Count += band.Count
				}
			}
		}
		for _, point := range item.Trend {
			counts[point.Label] += point.Count
			totals[point.Label] += point.Average * float64(point.Count)
		}
		merged.PassRate += item.PassRate
	}
	labels := make([]string, 0, len(counts))
	for label := range counts {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	for _, label := range labels {
		average := 0.0
		if counts[label] > 0 {
			average = RoundScore(totals[label] / float64(counts[label]))
		}
		merged.Trend = append(merged.Trend, TrendPoint{Label: label, Average: average, Count: counts[label]})
	}
	merged.PassRate = RoundScore(merged.PassRate / float64(len(items)))
	return merged
}

func Thresholds() map[string]float64 {
	return map[string]float64{"EXCELLENT": 90, "QUALIFIED": 75, "DEVELOPING": 60, "RETRAIN": 0}
}

func BandForScore(score float64) ScoreBand {
	for _, band := range BuildBands(nil) {
		if score >= band.Min && score <= band.Max {
			return band
		}
	}
	return ScoreBand{Name: "UNKNOWN"}
}
