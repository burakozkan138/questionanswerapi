package user_test

import (
	"burakozkan138/questionanswerapi/config"
	"burakozkan138/questionanswerapi/internal/api"
	"burakozkan138/questionanswerapi/internal/database"
	"burakozkan138/questionanswerapi/internal/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfile(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Run("Test Profile", func(t *testing.T) {
		if err := initDB(); err != nil {
			t.Error("Failed to initialize database")
		}
		routes := api.InitializeRoutes()

		req, err := http.NewRequest("GET", "/api/v1/auth/profile", nil)
		if err != nil {
			t.Error("Failed to create profile request")
		}

		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJxdWVzdGlvbmFuc3dlcmFwaSIsImV4cCI6MTcxMjUzMjEyNywiaWF0IjoxNzEyNDQ1NzI3LCJpc3MiOiJxdWVzdGlvbmFuc3dlcmFwaSIsInVzZXJJZCI6IjEwYTlhOTI3LWUzZTItNGQ5Yi1hODQ4LTUxZDY4OTY1NDUyMCJ9.9sNm7bHqxgB4Y5-lUQSRVXHDof52kElxyu6jYhNNyac")

		resp := httptest.NewRecorder()
		routes.ServeHTTP(resp, req)

		var baseResp models.Response
		err = json.Unmarshal(resp.Body.Bytes(), &baseResp)
		if err != nil {
			t.Error("Failed to unmarshal response body to Response")
		}

		if baseResp.HttpStatus != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, baseResp.HttpStatus)
		}

		if !baseResp.Success {
			t.Errorf("Expected success to be true, but got false")
		}

		if baseResp.Data == nil {
			t.Error("Expected data to be non-nil, but got nil")
		}

		t.Cleanup(func() {
			api.Rewrite()
		})
	})

	t.Run("Test Profile with invalid token", func(t *testing.T) {
		if err := initDB(); err != nil {
			t.Error("Failed to initialize database")
		}
		routes := api.InitializeRoutes()

		req, err := http.NewRequest("GET", "/api/v1/auth/profile", nil)
		if err != nil {
			t.Error("Failed to create profile request")
		}

		req.Header.Set("Authorization", "Bearer invalidtoken")

		resp := httptest.NewRecorder()
		routes.ServeHTTP(resp, req)

		var baseResp models.Response
		err = json.Unmarshal(resp.Body.Bytes(), &baseResp)
		if err != nil {
			t.Error("Failed to unmarshal response body to Response")
		}

		if baseResp.HttpStatus != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, baseResp.HttpStatus)
		}

		if baseResp.Success {
			t.Errorf("Expected success to be false, but got true")
		}

		if baseResp.Data != nil {
			t.Error("Expected data to be nil, but got non-nil")
		}

		t.Cleanup(func() {
			api.Rewrite()
		})
	})
}

func initDB() error {
	config.LoadConfig(".test", "env", "../../")
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
