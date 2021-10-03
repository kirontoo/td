package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

const dbName = "td.db"

var db *bolt.DB
var taskBucket = []byte("tasks")

type Task struct {
	Key   int
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

func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})

	if err != nil {
		return -1, err
	}

	// return id, nil
	return id, nil
}

// Integer to Byte
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Byte to Integer
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func Close() {
	db.Close()
}

func GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// only check keys because they hsould never be nil
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(key int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		id := itob(key)
		return tx.Bucket(taskBucket).Delete(id)
	})

	if err != nil {
		return err
	}

	return nil
}
