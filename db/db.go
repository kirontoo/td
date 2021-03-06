package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

const dbName = "td.db"

var db *bolt.DB
var taskBucket = []byte("tasks")

type Task struct {
	Key       int
	Value     string
	Completed bool
	Created   time.Time
}

func (t *Task) String() string {
	return fmt.Sprintf("%s", t.Value)
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

// Returns the key of the task
func CreateTask(taskValue string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)

		task := &Task{
			Key:       id,
			Value:     taskValue,
			Created:   time.Now(),
			Completed: false,
		}

		encoded, _ := json.Marshal(task)
		return b.Put(key, encoded)
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

// TODO: function to filter out completed tasks?
func GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t, err := unmarshalTask(v)
			if err != nil {
				return err
			}

			tasks = append(tasks, t)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetUncompletedTasks() ([]Task, error) {
	var tasks []Task

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t, err := unmarshalTask(v)
			if err != nil {
				return err
			}

			if !t.Completed {
				tasks = append(tasks, t)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetCompletedTasks() ([]Task, error) {
	var tasks []Task

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t, err := unmarshalTask(v)
			if err != nil {
				return err
			}

			if t.Completed {
				tasks = append(tasks, t)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(taskBucket).Delete(itob(key))
	})
}

func MarkCompleted(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		old := bucket.Get(itob(key))
		if old == nil {
			return errors.New("this task does not exist")
		}

		t, _ := unmarshalTask(old)

		// Don't bother marking completed if already completed
		if !t.Completed {
			t.Completed = true
		} else {
			return nil
		}

		encoded, _ := json.Marshal(t)
		return tx.Bucket(taskBucket).Put(itob(key), encoded)
	})
}

func unmarshalTask(t []byte) (Task, error) {
	var task Task
	err := json.Unmarshal(t, &task)
	if err != nil {
		return task, err
	}

	return task, nil
}
