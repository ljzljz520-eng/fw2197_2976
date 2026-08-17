package storage

import (
	"encoding/json"
	"fmt"
	"io"

	"go.etcd.io/bbolt"

	"skillsassessment/domain"
)

type RestoreReport struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}

func (d *Database) ImportJSON(reader io.Reader) (RestoreReport, error) {
	if reader == nil {
		return RestoreReport{}, fmt.Errorf("reader is required")
	}
	var projects []domain.AssessmentProject
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&projects); err != nil {
		return RestoreReport{}, err
	}
	report := RestoreReport{Errors: make([]string, 0)}
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return report, fmt.Errorf("database is closed")
	}
	err := d.db.Update(func(tx *bbolt.Tx) error {
		for _, project := range projects {
			if err := project.Validate(); err != nil {
				report.Skipped++
				report.Errors = append(report.Errors, project.ID+": "+err.Error())
				continue
			}
			data, err := json.Marshal(project)
			if err != nil {
				return err
			}
			if err := tx.Bucket(bucketProjects).Put([]byte(project.ID), data); err != nil {
				return err
			}
			report.Imported++
		}
		return nil
	})
	return report, err
}

func (d *Database) ForEachProject(fn func(domain.AssessmentProject) error) error {
	if fn == nil {
		return fmt.Errorf("callback is required")
	}
	return d.ReadTransaction(func(tx *Transaction) error {
		return tx.eachProject(fn)
	})
}

func (t *Transaction) eachProject(fn func(domain.AssessmentProject) error) error {
	return t.tx.Bucket(bucketProjects).ForEach(func(_, value []byte) error {
		var project domain.AssessmentProject
		if err := json.Unmarshal(value, &project); err != nil {
			return err
		}
		return fn(project)
	})
}
