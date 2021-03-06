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

var (
	ErrTaskDoesNotExist = errors.New("Task does not exist")
	ErrTaskNotCreated   = errors.New("Task not created")
	ErrNoTasks          = errors.New("No tasks returned")
	ErrNotEnoughTasks   = errors.New("Not enough tasks returned")
)

func TestInitDb(t *testing.T) {
	t.Run("should initialize a new database", func(t *testing.T) {
		dbPath := getTestDbPath()
		err := openDb()
		if err != nil {
			t.FailNow()
		}

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
		deleteDb()
	})
}

func TestCreateTask(t *testing.T) {
	openDb()

	t.Run("should create a new task", func(t *testing.T) {
		want := "this is a new task"
		key, err := CreateTask(want)

		if err != nil {
			t.Error(ErrTaskNotCreated)
			t.Failed()
		}

		var got []byte
		// check task was created
		if err := db.View(func(tx *bolt.Tx) error {
			got = tx.Bucket(taskBucket).Get(itob(key))
			if got == nil {
				return ErrTaskDoesNotExist
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
			t.Error(ErrTaskNotCreated)
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

			if got.Value != want && got.Key != key {
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

func TestMarkCompleted(t *testing.T) {
	openDb()

	t.Run("should mark a task as completed", func(t *testing.T) {
		taskStr := "Test mark completed"
		key, cErr := CreateTask(taskStr)
		if cErr != nil {
			t.Error(ErrTaskNotCreated)
			t.FailNow()
		}

		// Task.Completed should be true
		want := true

		// Mark completed
		err := MarkCompleted(key)
		if err != nil {
			t.Error(err)
		}

		var got []byte
		if err := db.View(func(tx *bolt.Tx) error {
			got = tx.Bucket(taskBucket).Get(itob(key))
			if got == nil {
				return ErrTaskDoesNotExist
			}
			return nil
		}); err != nil {
			t.Log(err)
			t.FailNow()
		}

		task, _ := unmarshalTask(got)
		if task.Completed != want {
			t.Log(task)
			t.Error("Task was not marked completed")
			t.Fail()
		}

	})

	t.Run("fail to mark complete task that doesn't exist", func(t *testing.T) {
		err := MarkCompleted(80)
		if err == nil {
			t.Error("Marked a task that should not exist")
			t.Fail()
		}
	})

	t.Cleanup(func() {
		db.Close()
		deleteDb()
	})
}

func TestGetTasks(t *testing.T) {
	openDb()
	orgTasks := []string{"task 1", "task 2", "task 3", "task 4"}

	t.Run("Get all tasks", func(t *testing.T) {
		var keys []int

		for _, t := range orgTasks {
			k, _ := CreateTask(t)
			keys = append(keys, k)
		}

		tasks, err := GetAllTasks()
		if err != nil {
			t.Error(ErrNoTasks)
			t.Fail()
		}

		if len(keys) != len(tasks) {
			t.Error(ErrNotEnoughTasks)
			t.Fail()
		}
	})

	t.Run("Get all uncompleted tasks", func(t *testing.T) {
		tasks, err := GetUncompletedTasks()
		if err != nil {
			t.Error(ErrNoTasks)
			t.Fail()
		}

		if len(orgTasks) != len(tasks) {
			t.Error(ErrNotEnoughTasks)
			t.Fail()
		}

		for _, task := range tasks {
			if task.Completed {
				t.Error("Returned a completed task")
				t.Fail()
			}
		}
	})

	t.Run("Get all completed tasks", func(t *testing.T) {
		MarkCompleted(1)
		MarkCompleted(2)

		tasks, err := GetCompletedTasks()

		if err != nil {
			t.Error(ErrNoTasks)
			t.Fail()
		}

		if len(tasks) != 2 {
			t.Error(ErrNotEnoughTasks)
			t.Fail()
		}

		for _, task := range tasks {
			if !task.Completed {
				t.Error("Returned a completed task")
				t.Fail()
			}
		}
	})

	t.Cleanup(func() {
		db.Close()
		deleteDb()
	})
}

func TestDeleteTask(t *testing.T) {
	openDb()

	key, err := CreateTask("hello world")
	if err != nil {
		t.Error(ErrTaskNotCreated)
		t.FailNow()
	}

	err = DeleteTask(key)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	tasks, _ := GetAllTasks()
	if len(tasks) != 0 {
		t.Error("The task was not deleted")
		t.Fail()
	}

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
