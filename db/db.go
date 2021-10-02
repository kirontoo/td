package db

import (
	"time"

	"github.com/boltdb/bolt"
)

const dbName = "td.db"

var db *bolt.DB
var taskBucket = []byte("tasks")

type Task struct {
	ID    int
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func AddTask(task Task) error {
	return nil
}
