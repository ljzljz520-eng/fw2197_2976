package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"go.etcd.io/bbolt"

	"skillsassessment/domain"
)

type Query struct {
	Prefix string
	Status domain.ProjectStatus
}

func (d *Database) QueryProjects(query Query) ([]domain.AssessmentProject, error) {
	projects, err := d.ListProjects()
	if err != nil {
		return nil, err
	}
	selected := make([]domain.AssessmentProject, 0, len(projects))
	for _, project := range projects {
		if query.Prefix != "" && !strings.HasPrefix(project.ID, query.Prefix) {
			continue
		}
		if query.Status != "" && project.Status != query.Status {
			continue
		}
		selected = append(selected, project)
	}
	return selected, nil
}

func (d *Database) Count(bucketName string) (int, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return 0, fmt.Errorf("database is closed")
	}
	count := 0
	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}
		count = bucket.Stats().KeyN
		return nil
	})
	return count, err
}

func (d *Database) ExportJSON() ([]byte, error) {
	projects, err := d.ListProjects()
	if err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Database) HasProject(id string) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return false, fmt.Errorf("database is closed")
	}
	var found bool
	err := d.db.View(func(tx *bbolt.Tx) error {
		found = tx.Bucket(bucketProjects).Get([]byte(id)) != nil
		return nil
	})
	return found, err
}

func EqualJSON(left, right []byte) bool {
	return bytes.Equal(bytes.TrimSpace(left), bytes.TrimSpace(right))
}
