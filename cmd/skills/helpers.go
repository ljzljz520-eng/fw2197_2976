package main

import (
	"errors"
	"fmt"

	"skillsassessment/domain"
)

func parseCertification(kind, value string) (domain.Certification, error) {
	switch kind {
	case "PASS", "FAIL":
		return domain.Certification{Kind: domain.CertificationGrade, Grade: kind}, nil
	case "SCORE":
		return domain.Certification{Kind: domain.CertificationScore, Percent: 100}, nil
	default:
		return domain.Certification{}, errors.New("unsupported certification " + fmt.Sprint(value))
	}
}
