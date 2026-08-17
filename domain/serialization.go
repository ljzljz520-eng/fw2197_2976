package domain

import (
	"encoding/json"
	"errors"
	"strings"
)

type ProjectEnvelope struct {
	Schema  int               `json:"schema"`
	Project AssessmentProject `json:"project"`
}

func EncodeProject(project AssessmentProject) ([]byte, error) {
	if err := project.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(ProjectEnvelope{Schema: 1, Project: project})
}

func DecodeProject(data []byte) (AssessmentProject, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return AssessmentProject{}, errors.New("project payload is empty")
	}
	var envelope ProjectEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return AssessmentProject{}, err
	}
	if envelope.Schema != 1 {
		return AssessmentProject{}, errors.New("unsupported project schema")
	}
	if err := envelope.Project.Validate(); err != nil {
		return AssessmentProject{}, err
	}
	return envelope.Project, nil
}

func EncodeLearner(learner Learner) ([]byte, error) {
	if err := learner.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(learner)
}

func DecodeLearner(data []byte) (Learner, error) {
	var learner Learner
	if err := json.Unmarshal(data, &learner); err != nil {
		return learner, err
	}
	if err := learner.Validate(); err != nil {
		return learner, err
	}
	return learner, nil
}

func ProjectFingerprint(project AssessmentProject) string {
	return project.ID + ":" + string(project.Status) + ":" + itoa(project.Version) + ":" + itoa(len(project.Learners))
}
