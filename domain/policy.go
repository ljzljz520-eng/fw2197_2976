package domain

import "strings"

type ReviewPolicy struct {
	MinimumLearners int
	MinimumAverage  float64
	RequireSafety   bool
}

type PolicyDecision struct {
	Allowed bool
	Reasons []string
}

func (p ReviewPolicy) Evaluate(project AssessmentProject) PolicyDecision {
	reasons := make([]string, 0)
	if len(project.Learners) < p.MinimumLearners {
		reasons = append(reasons, "not enough learners")
	}
	summary := SummarizeScores(project.Learners)
	if p.MinimumAverage > 0 && summary.Average < p.MinimumAverage {
		reasons = append(reasons, "average score is below policy")
	}
	if p.RequireSafety {
		for _, learner := range project.Learners {
			if learner.Safety < p.MinimumAverage {
				reasons = append(reasons, "safety score is below policy for "+learner.ID)
			}
		}
	}
	return PolicyDecision{Allowed: len(reasons) == 0, Reasons: reasons}
}

func DefaultReviewPolicy() ReviewPolicy {
	return ReviewPolicy{MinimumLearners: 1, MinimumAverage: 60, RequireSafety: true}
}

func NormalizeGender(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	switch value {
	case "M", "MALE":
		return "M"
	case "F", "FEMALE":
		return "F"
	default:
		return "U"
	}
}

func NormalizeCertification(c Certification) Certification {
	if c.Kind == CertificationGrade {
		c.Grade = strings.ToUpper(strings.TrimSpace(c.Grade))
		return c
	}
	if c.Percent < 0 {
		c.Percent = 0
	}
	if c.Percent > 100 {
		c.Percent = 100
	}
	return c
}
