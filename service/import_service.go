package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"skillsassessment/domain"
)

func (s *ProjectService) ImportLearners(projectID string, input io.Reader) (int, error) {
	if input == nil {
		return 0, errors.New("input is required")
	}
	scanner := bufio.NewScanner(input)
	line := 0
	added := 0
	for scanner.Scan() {
		line++
		fields := strings.Split(scanner.Text(), ",")
		if len(fields) != 7 {
			return added, fmt.Errorf("line %d: expected seven fields", line)
		}
		theory, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return added, fmt.Errorf("line %d theory: %w", line, err)
		}
		practical, err := strconv.ParseFloat(fields[5], 64)
		if err != nil {
			return added, fmt.Errorf("line %d practical: %w", line, err)
		}
		safety, err := strconv.ParseFloat(fields[6], 64)
		if err != nil {
			return added, fmt.Errorf("line %d safety: %w", line, err)
		}
		learner := domain.Learner{ID: fields[0], Name: fields[1], Gender: fields[2], Email: fields[3], Theory: theory, Practical: practical, Safety: safety}
		if err := s.AddLearner(projectID, learner); err != nil {
			return added, fmt.Errorf("line %d: %w", line, err)
		}
		added++
	}
	if err := scanner.Err(); err != nil {
		return added, err
	}
	return added, nil
}
