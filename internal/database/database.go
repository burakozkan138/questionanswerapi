package database

import (
	"burakozkan138/questionanswerapi/config"
	"burakozkan138/questionanswerapi/internal/models"
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once     = sync.Once{}
	db       *gorm.DB
	dbConfig *config.DatabaseConfig
)

func GetDB() *gorm.DB {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.Port)

		newDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})

		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Connected to database %s successfully", newDB.Migrator().CurrentDatabase())

		db = newDB
	})
	return db
}

func InitializeDatabase(cfg *config.DatabaseConfig) error {
	dbConfig = cfg
	db := GetDB()

	if err := db.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.QuestionLike{},
		&models.Answer{},
		&models.AnswerLike{},
	); err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Database migrated successfully")
	return nil
}

func Close() {
	db, _ := db.DB()
	db.Close()
	log.Printf("Database connection closed")
}
