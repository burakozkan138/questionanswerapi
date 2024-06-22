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
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		requestUser models.RegisterValidation
		dbUser      models.User
		db          = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		response := models.NewResponse(false, "Invalid JSON", http.StatusBadRequest, nil, nil)
		response.Write(w)
		return
	}

	if err := db.Where("email = ?", requestUser.Email).Or("username = ?", requestUser.Username).First(&dbUser).Error; err == nil {
		response := models.NewResponse(false, "User already exists", http.StatusConflict, nil, nil)
		response.Write(w)
		return
	}

	dbUser = models.User{
		Fullname: requestUser.Fullname,
		Username: requestUser.Username,
		Email:    requestUser.Email,
		Password: requestUser.Password,
	}

	if err := db.Create(&dbUser).Error; err != nil {
		response := models.NewResponse(false, "User could not be created", http.StatusInternalServerError, nil, nil)
		response.Write(w)
		return
	}

	authResponse, err := dbUser.ToUserAuthResponse()
	if err != nil {
		response := models.NewResponse(false, "Internal server error", http.StatusInternalServerError, nil, nil)
		response.Write(w)
		return
	}

	response := models.NewResponse(true, "User created successfully", http.StatusCreated, authResponse, nil)
	response.Write(w)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email or username the provided information
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginValidation true "User information"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		requestUser models.LoginValidation
		dbUser      models.User
		db          = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		response := models.NewResponse(false, "Invalid JSON", http.StatusBadRequest, nil, nil)
		response.Write(w)
		return
	}

	if err := db.Where("email = ?", requestUser.Email).Or("username = ?", requestUser.Username).First(&dbUser).Error; err != nil {
		response := models.NewResponse(false, "User not found", http.StatusNotFound, nil, nil)
		response.Write(w)
		return
	}

	if !dbUser.CheckPassword(requestUser.Password) {
		repsonse := models.NewResponse(false, "Invalid password", http.StatusUnauthorized, nil, nil)
		repsonse.Write(w)
		return
	}

	if dbUser.Blocked {
		response := models.NewResponse(false, "User is blocked", http.StatusForbidden, nil, nil)
		response.Write(w)
		return
	}

	authResponse, err := dbUser.ToUserAuthResponse()
	if err != nil {
		response := models.NewResponse(false, "Internal server error", http.StatusInternalServerError, nil, nil)
		response.Write(w)
		return
	}

	response := models.NewResponse(true, "User logged in successfully", http.StatusOK, authResponse, nil)
	response.Write(w)
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
		response := models.NewResponse(false, "Invalid JSON", http.StatusBadRequest, nil, nil)
		response.Write(w)
		return
	}

	if err := db.Where("email = ?", requestEmail.Email).First(&dbUser).Error; err != nil {
		response := models.NewResponse(false, "User not found", http.StatusNotFound, nil, nil)
		response.Write(w)
		return
	}

	// TODO: Tamamla
	// Send email
	// pkg.SendEmail()

	w.Write([]byte("Forgot password"))
}
