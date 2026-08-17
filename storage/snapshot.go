package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Snapshot struct {
	CreatedAt string `json:"created_at"`
	Projects  int    `json:"projects"`
	Learners  int    `json:"learners"`
	Audits    int    `json:"audits"`
}

func (d *Database) Snapshot() (Snapshot, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return Snapshot{}, fmt.Errorf("database is closed")
	}
	result := Snapshot{CreatedAt: time.Unix(0, 0).UTC().Format(time.RFC3339)}
	err := d.db.View(func(tx *bbolt.Tx) error {
		result.Projects = tx.Bucket(bucketProjects).Stats().KeyN
		result.Learners = tx.Bucket(bucketLearners).Stats().KeyN
		result.Audits = tx.Bucket(bucketAudits).Stats().KeyN
		return nil
	})
	return result, err
}

func (s Snapshot) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
