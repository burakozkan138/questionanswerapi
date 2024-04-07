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

func TestRegister(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Run("Test Register", func(t *testing.T) {
		if err := initDB(); err != nil {
			t.Error("Failed to initialize database")
		}
		routes := api.InitilazeRoutes()

		req, err := createRegisterRequest("Test User", "testuser", "test@test.com", "testuser", "testuser")
		if err != nil {
			t.Error("Failed to create register request")
		}

		resp := httptest.NewRecorder()
		routes.ServeHTTP(resp, req)

		var baseResp models.BaseResponse
		err = json.Unmarshal(resp.Body.Bytes(), &baseResp)
		if err != nil {
			t.Error("Failed to unmarshal response body to BaseResponse")
		}

		if baseResp.HttpStatus != http.StatusCreated {
			t.Errorf("Expected status code %d, but got %d", http.StatusCreated, baseResp.HttpStatus)
		}

		if baseResp.Message != "User created successfully" {
			t.Errorf("Expected message %s, but got %s", "User created successfully", baseResp.Message)
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

	t.Run("Test Register with invalid data", func(t *testing.T) {
		if err := initDB(); err != nil {
			t.Error("Failed to initialize database")
		}

		routes := api.InitilazeRoutes()

		req, err := createRegisterRequest("", "", "invalidemail", "", "")
		if err != nil {
			t.Error("Failed to create register request")
		}

		resp := httptest.NewRecorder()
		routes.ServeHTTP(resp, req)

		var errResp models.ErrorResponse
		err = json.Unmarshal(resp.Body.Bytes(), &errResp)
		if err != nil {
			t.Error("Failed to unmarshal response body to BaseResponse")
		}

		if errResp.HttpStatus != http.StatusBadRequest {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, errResp.HttpStatus)
		}

		if errResp.Message != "Validation error" {
			t.Errorf("Expected message %s, but got %s", "Validation error", errResp.Message)
		}

		if errResp.Success {
			t.Errorf("Expected success to be false, but got true")
		}

		if errResp.Errors == nil {
			t.Error("Expected errors to be non-nil, but got nil")
		}

	})
}

func createRegisterRequest(fullname, username, email, password, confirmPassword string) (*http.Request, error) {
	reqBody := models.RegisterValidation{
		Fullname:        fullname,
		Username:        username,
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	reqBodyReader := bytes.NewBuffer(reqBodyJSON)

	req, err := http.NewRequest("POST", "/api/v1/auth/register", reqBodyReader)
	if err != nil {
		return nil, err
	}

	return req, nil
}
