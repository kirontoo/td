package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

const testDbName = "test.db"

func TestInitDb(t *testing.T) {
	// TODO create a after test to delete the created db

	t.Run("should initialize a new database", func(t *testing.T) {
		dbPath := getTestDbPath()
		err := openDb()
		if err != nil {
			t.FailNow()
		}

		// TODO: create a new test for this
		// created database in wrong path
		currDbPath := db.Path()
		if currDbPath != dbPath {
			t.Errorf("DB created in wrong path.\ngot: %s\n want: %s", currDbPath, dbPath)
		}

		_, fileErr := os.Stat(dbPath)
		if os.IsNotExist(fileErr) {
			t.Error("Expected database to be created")
		}
	})

	t.Cleanup(func() {
		db.Close()
	})
}

func TestCreateTask(t *testing.T) {
	openDb()

	t.Run("should create a new task", func(t *testing.T) {
		want := "this is a new task"
		key, err := CreateTask(want)

		if err != nil {
			t.Errorf("should have created a new task")
			t.Failed()
		}

		var got []byte
		// check task was created
		if err := db.View(func(tx *bolt.Tx) error {
			got = tx.Bucket(taskBucket).Get(itob(key))
			if got == nil {
				return errors.New("this task does not exist")
			}
			return nil
		}); err != nil {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("should have created a task with correct value", func(t *testing.T) {
		want := "this is a new task"
		key, err := CreateTask(want)

		if err != nil {
			t.Errorf("should have created a new task")
			t.Failed()
		}

		if err := db.View(func(tx *bolt.Tx) error {
			var got Task
			value := tx.Bucket(taskBucket).Get(itob(key))
			jsonErr := json.Unmarshal([]byte(value), &got)
			if jsonErr != nil {
				t.Error("Failed to unmarshal value")
				t.FailNow()
			}

			if got.Value != want {
				return fmt.Errorf("Created the wrong task.\n got: %s\n want: %s", got.Value, want)
			}

			return nil
		}); err != nil {
			t.Log(err)
			t.Fail()
		}
	})

	t.Cleanup(func() {
		db.Close()
		deleteDb()
	})
}

func getTestDbPath() string {
	currDir, _ := os.Getwd()
	return filepath.Join(path.Dir(currDir), testDbName)
}

func openDb() error {
	dbPath := getTestDbPath()
	return Init(dbPath)
}

func deleteDb() error {
	pathToFile := getTestDbPath()
	return os.Remove(pathToFile)
}
