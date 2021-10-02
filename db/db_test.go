package db

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestInitDb(t *testing.T) {
	// TODO create a after test to delete the created db

	t.Run("should initialize a new database", func(t *testing.T) {
		currDir, _ := os.Getwd()
		dbPath := filepath.Join(path.Dir(currDir), "test.db")
		Init(dbPath)

		_, fileErr := os.Stat(dbPath)
		if os.IsNotExist(fileErr) {
			t.Error("Expected database to be created")
		}
	})
}
