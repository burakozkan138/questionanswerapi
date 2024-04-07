package auth_test

import (
	"burakozkan138/questionanswerapi/config"
	"burakozkan138/questionanswerapi/internal/database"
)

func initDB() error {
	err := config.LoadConfig(".env.test", "env", "../../config")
	if err != nil {
		return err
	}
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	err = database.InitializeDatabase(&cfg.Database)
	if err != nil {
		return err
	}

	return nil
}
