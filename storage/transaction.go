package storage

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

type Transaction struct {
	db *Database
	tx *bbolt.Tx
}

func (d *Database) ReadTransaction(fn func(*Transaction) error) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return errors.New("database is closed")
	}
	return d.db.View(func(tx *bbolt.Tx) error {
		return fn(&Transaction{db: d, tx: tx})
	})
}

func (t *Transaction) ProjectIDs() ([]string, error) {
	ids := make([]string, 0)
	if t.tx == nil {
		return nil, fmt.Errorf("transaction is closed")
	}
	err := t.tx.Bucket(bucketProjects).ForEach(func(key, _ []byte) error {
		ids = append(ids, string(key))
		return nil
	})
	return ids, err
}
