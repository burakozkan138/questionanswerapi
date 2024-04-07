package database_test

import (
	"burakozkan138/questionanswerapi/config"
	"burakozkan138/questionanswerapi/internal/database"
	"io"
	"log"
	"testing"

	"gorm.io/gorm"
)

func TestDatabase(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Run("Test Initialize Database", func(t *testing.T) {
		_, err := InitDB()

		if err != nil {
			t.Error("Database initialization failed")
		}
	})

	t.Run("Test Get Database", func(t *testing.T) {
		db, err := InitDB()
		if err != nil {
			t.Error("Database initialization failed")
		}

		if db == nil {
			t.Error("Database is nil")
		}
	})

	t.Run("Test Ping Database", func(t *testing.T) {
		db, err := InitDB()
		if err != nil {
			t.Error("Database initialization failed")
		}

		if err := db.Exec("SELECT 1").Error; err != nil {
			t.Error("Database connection failed")
		}
	})

	t.Cleanup(func() {
		database.Close()
	})
}

func InitDB() (*gorm.DB, error) {
	err := config.LoadConfig(".env.test", "env", "../../config")
	cfg, err := config.InitializeConfig()
	if err != nil {
		return nil, err
	}

	err = database.InitializeDatabase(&cfg.Database)
	if err != nil {
		return nil, err
	}

	return database.GetDB(), nil
}
