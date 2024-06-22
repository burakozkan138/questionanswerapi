package main

import (
	"burakozkan138/questionanswerapi/config"
	"burakozkan138/questionanswerapi/internal/api"
	"burakozkan138/questionanswerapi/internal/database"
	"fmt"
	"log"
	"net/http"
)

// @title Question Answer API
// @description This API provides endpoints for managing questions and answers.
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := config.LoadConfig(".env", "env", "../config")
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
		return
	}

	config, err := config.InitializeConfig()
	if err != nil {
		return
	}

	err = database.InitializeDatabase(&config.Database)
	if err != nil {
		return
	}

	handler := api.InitializeRoutes()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: handler,
	}

	log.Printf("Server listening on port :%d\n", config.Server.Port)
	server.ListenAndServe()
}
