package users

import (
	"burakozkan138/questionanswerapi/internal/database"
	"burakozkan138/questionanswerapi/internal/models"
	"encoding/json"
	"net/http"
)

type Handler struct{}

// Fetch profile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/profile [get]
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	var (
		dbUser models.User
		db     = database.GetDB()
	)

	userID := r.Context().Value(models.USER_CTX_KEY)

	if err := db.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		response := models.NewResponse(false, "User not found", http.StatusNotFound, nil, nil)
		response.Write(w)
		return
	}

	response := models.NewResponse(true, "User fetched successfully", http.StatusOK, dbUser.ToUserProfileResponse(), nil)
	response.Write(w)
}

// Edit godoc
// @Summary Edit user profile
// @Description Edit the profile of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body models.EditUserValidation true "User information"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/edit [put]
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	var (
		requestUser models.EditUserValidation
		dbUser      models.User
		db          = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		response := models.NewResponse(false, "Invalid JSON", http.StatusBadRequest, nil, nil)
		response.Write(w)
		return
	}

	userID := r.Context().Value(models.USER_CTX_KEY)

	if err := db.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		response := models.NewResponse(false, "User not found", http.StatusNotFound, nil, nil)
		response.Write(w)
		return
	}

	dbUser.EditUserCheckFields(requestUser)

	if err := db.Save(&dbUser).Error; err != nil {
		response := models.NewResponse(false, "User edit failed", http.StatusInternalServerError, nil, nil)
		response.Write(w)
		return
	}

	response := models.NewResponse(true, "User edit successfully", http.StatusOK, dbUser.ToUserProfileResponse(), nil)
	response.Write(w)
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload")) // todo: implement image upload
}
