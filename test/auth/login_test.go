package auth_test

import (
	"burakozkan138/questionanswerapi/internal/api"
	"burakozkan138/questionanswerapi/internal/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	log.SetOutput(io.Discard)

	t.Run("Test Login", func(t *testing.T) {
		if err := initDB(); err != nil {
			t.Error("Failed to initialize database")
		}
		routes := api.InitilazeRoutes()

		req, err := createLoginRequest("testuser", "testuser")
		if err != nil {
			t.Error("Failed to create login request")
		}

		resp := httptest.NewRecorder()
		routes.ServeHTTP(resp, req)

		var response models.BaseResponse
		err = json.Unmarshal(resp.Body.Bytes(), &response)
		if err != nil {
			t.Error("Failed to unmarshal response body to BaseResponse")
		}

		if response.HttpStatus != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.HttpStatus)
		}

		if !response.Success {
			t.Errorf("Expected success to be true, but got false")
		}

		if response.Data == nil {
			t.Error("Expected data to be not nil")
		}

		t.Cleanup(func() {
			api.Rewrite()
		})
	})

	t.Run("Test Login with invalid credentials", func(t *testing.T) {
		if err := initDB(); err != nil {
			t.Error("Failed to initialize database")
		}
		routes := api.InitilazeRoutes()

		req, err := createLoginRequest("testuser", "invalidpassword")
		if err != nil {
			t.Error("Failed to create login request")
		}

		resp := httptest.NewRecorder()
		routes.ServeHTTP(resp, req)

		var response models.BaseResponse
		err = json.Unmarshal(resp.Body.Bytes(), &response)
		if err != nil {
			t.Error("Failed to unmarshal response body to BaseResponse")
		}

		if response.HttpStatus != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, response.HttpStatus)
		}

		if response.Success {
			t.Errorf("Expected success to be false, but got true")
		}

		if response.Data != nil {
			t.Error("Expected data to be nil")
		}

		t.Cleanup(func() {
			api.Rewrite()
		})
	})
}

func createLoginRequest(username, password string) (*http.Request, error) {
	reqBody := models.LoginValidation{
		Username: username,
		Password: password,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	reqBodyReader := bytes.NewBuffer(reqBodyJSON)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", reqBodyReader)
	if err != nil {
		return nil, err
	}

	return req, nil
}
