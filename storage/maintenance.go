package storage

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type MaintenanceResult struct {
	Buckets []string `json:"buckets"`
	Keys    int      `json:"keys"`
	Checked string   `json:"checked"`
}

func (d *Database) CheckIntegrity() (MaintenanceResult, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return MaintenanceResult{}, fmt.Errorf("database is closed")
	}
	result := MaintenanceResult{Buckets: []string{"assessment_projects", "learners", "certifications", "audit_trails"}, Checked: time.Unix(0, 0).UTC().Format(time.RFC3339)}
	err := d.db.View(func(tx *bbolt.Tx) error {
		for _, name := range result.Buckets {
			bucket := tx.Bucket([]byte(name))
			if bucket == nil {
				return fmt.Errorf("missing bucket %s", name)
			}
			result.Keys += bucket.Stats().KeyN
		}
		return nil
	})
	return result, err
}

func (d *Database) CompactTo(target string) error {
	if target == "" {
		return fmt.Errorf("target path is required")
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return fmt.Errorf("database is closed")
	}
	return bbolt.Compact(nil, d.db, 0)
}

func (d *Database) UpdateMetadata(key, value string) error {
	if key == "" {
		return fmt.Errorf("metadata key is required")
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return fmt.Errorf("database is closed")
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("metadata"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (d *Database) ReadMetadata(key string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return "", fmt.Errorf("database is closed")
	}
	var value []byte
	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("metadata"))
		if bucket == nil {
			return ErrNotFound
		}
		value = append([]byte(nil), bucket.Get([]byte(key))...)
		if value == nil {
			return ErrNotFound
		}
		return nil
	})
	return string(value), err
}

func (d *Database) DeleteMetadata(key string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return fmt.Errorf("database is closed")
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("metadata"))
		if bucket == nil {
			return ErrNotFound
		}
		if bucket.Get([]byte(key)) == nil {
			return ErrNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func (d *Database) FlushMetadata() error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return fmt.Errorf("database is closed")
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		if tx.Bucket([]byte("metadata")) == nil {
			return nil
		}
		return tx.DeleteBucket([]byte("metadata"))
	})
}
