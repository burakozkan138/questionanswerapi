// Package auth provides authentication handlers.
package auth

import (
	"burakozkan138/questionanswerapi/internal/database"
	"burakozkan138/questionanswerapi/internal/models"
	"encoding/json"
	"net/http"
)

type Handler struct{}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.RegisterValidation true "User information"
// @Success 201 {object} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 409 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		requestUser models.RegisterValidation
		dbUser      models.User
		db          = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.Where("email = ?", requestUser.Email).Or("username = ?", requestUser.Username).First(&dbUser).Error; err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	dbUser = models.User{
		Fullname: requestUser.Fullname,
		Username: requestUser.Username,
		Email:    requestUser.Email,
		Password: requestUser.Password,
	}

	if err := db.Create(&dbUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResponse, err := dbUser.ToUserAuthResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "User created successfully", http.StatusCreated, authResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusCreated)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email or username the provided information
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginValidation true "User information"
// @Success 200 {object} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 401 {object} models.BaseResponse
// @Failure 403 {object} models.BaseResponse
// @Failure 404 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		requestUser models.LoginValidation
		dbUser      models.User
		db          = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.Where("email = ?", requestUser.Email).Or("username = ?", requestUser.Username).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !dbUser.CheckPassword(requestUser.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	if dbUser.Blocked {
		http.Error(w, "User is blocked", http.StatusForbidden)
		return
	}

	authResponse, err := dbUser.ToUserAuthResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "User logged in successfully", http.StatusOK, authResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Tamamla
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var (
		requestEmail models.ForgotPasswordValidation
		dbUser       models.User
		db           = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestEmail); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.Where("email = ?", requestEmail.Email).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// TODO: Tamamla
	// Send email
	// pkg.SendEmail()

	w.Write([]byte("Forgot password"))
}
