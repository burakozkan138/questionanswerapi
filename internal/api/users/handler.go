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
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /user/profile [get]
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	var (
		dbUser models.User
		db     = database.GetDB()
	)

	userID := r.Context().Value("userID")

	if err := db.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "User fetched successfully", http.StatusOK, dbUser.ToUserProfileResponse())

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Edit godoc
// @Summary Edit user profile
// @Description Edit the profile of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body models.EditUserValidation true "User information"
// @Success 200 {object} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /user/edit [put]
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	var (
		requestUser models.EditUserValidation
		dbUser      models.User
		db          = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")

	if err := db.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbUser.EditUserCheckFields(requestUser)

	if err := db.Save(&dbUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "User edit successfully", http.StatusOK, dbUser.ToUserProfileResponse())

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload")) // todo: implement image upload
}
