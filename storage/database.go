package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.etcd.io/bbolt"

	"skillsassessment/domain"
)

var (
	ErrNotFound          = errors.New("record not found")
	bucketProjects       = []byte("assessment_projects")
	bucketLearners       = []byte("learners")
	bucketCertifications = []byte("certifications")
	bucketAudits         = []byte("audit_trails")
)

type Database struct {
	path string
	db   *bbolt.DB
	mu   sync.RWMutex
}

func Open(path string) (*Database, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}
	bolt, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	database := &Database{path: path, db: bolt}
	if err := database.initialize(); err != nil {
		_ = bolt.Close()
		return nil, err
	}
	return database, nil
}

func (d *Database) initialize() error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{bucketProjects, bucketLearners, bucketCertifications, bucketAudits} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})
}

func (d *Database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.db == nil {
		return nil
	}
	err := d.db.Close()
	d.db = nil
	return err
}

func (d *Database) Path() string {
	return d.path
}

func (d *Database) PutProject(project domain.AssessmentProject) error {
	if err := project.Validate(); err != nil {
		return err
	}
	data, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("encode project: %w", err)
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return errors.New("database is closed")
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(bucketProjects).Put([]byte(project.ID), data); err != nil {
			return err
		}
		if err := putJSON(tx.Bucket(bucketAudits), project.ID, project.Audit); err != nil {
			return err
		}
		for _, learner := range project.Learners {
			if err := putJSON(tx.Bucket(bucketLearners), learner.ID, learner); err != nil {
				return err
			}
		}
		return putJSON(tx.Bucket(bucketCertifications), project.ID, project.Certification)
	})
}

func (d *Database) GetProject(id string) (domain.AssessmentProject, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var project domain.AssessmentProject
	if d.db == nil {
		return project, errors.New("database is closed")
	}
	err := d.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucketProjects).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return json.Unmarshal(append([]byte(nil), value...), &project)
	})
	if err != nil {
		return project, err
	}
	return project, nil
}

func (d *Database) DeleteProject(id string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return errors.New("database is closed")
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		if tx.Bucket(bucketProjects).Get([]byte(id)) == nil {
			return ErrNotFound
		}
		if err := tx.Bucket(bucketProjects).Delete([]byte(id)); err != nil {
			return err
		}
		return tx.Bucket(bucketCertifications).Delete([]byte(id))
	})
}

func (d *Database) ListProjects() ([]domain.AssessmentProject, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	projects := make([]domain.AssessmentProject, 0)
	if d.db == nil {
		return nil, errors.New("database is closed")
	}
	err := d.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketProjects).ForEach(func(_, value []byte) error {
			var project domain.AssessmentProject
			if err := json.Unmarshal(value, &project); err != nil {
				return err
			}
			projects = append(projects, project)
			return nil
		})
	})
	return projects, err
}

func putJSON(bucket *bbolt.Bucket, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode record %s: %w", key, err)
	}
	return bucket.Put([]byte(key), data)
}
