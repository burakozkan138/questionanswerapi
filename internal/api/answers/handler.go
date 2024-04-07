package answers

import (
	"burakozkan138/questionanswerapi/internal/database"
	"burakozkan138/questionanswerapi/internal/models"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct{}

// Create Answer godoc
// @Summary Create a new answer
// @Description Create a new answer for the provided question
// @Tags Answers
// @Accept json
// @Produce json
// @Security Bearer
// @Param question_id path string true "Question ID"
// @Param body body models.CreateAnswerValidation true "Answer information"
// @Success 201 {object} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /answer/{question_id}/answers [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		requestAnswer models.CreateAnswerValidation
		dbAnswer      models.Answer
		db            = database.GetDB()
	)

	questionID := r.PathValue("question_id")
	userID := r.Context().Value("userID")

	if err := json.NewDecoder(r.Body).Decode(&requestAnswer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbAnswer = models.Answer{
		Content:    requestAnswer.Content,
		QuestionID: uuid.MustParse(questionID),
		UserID:     userID.(uuid.UUID),
	}

	if err := db.Create(&dbAnswer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Preload("User").First(&dbAnswer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Answer created successfully", http.StatusCreated, dbAnswer.ToResponse())

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusCreated)
}

// Fetch answers godoc
// @Summary Get answers for a question
// @Description Get answers for the provided question
// @Tags Answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /answer/{question_id}/answers [get]
func (h *Handler) GetAnswers(w http.ResponseWriter, r *http.Request) {
	var (
		dbAnswers []models.Answer
		db        = database.GetDB()
	)

	questionID := r.PathValue("question_id")

	if err := db.Preload("User").Preload("Likes").Where("question_id = ?", questionID).Find(&dbAnswers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Answers fetched successfully", http.StatusOK, models.MapToAnswerResponse(dbAnswers))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Like Answer godoc
// @Summary Like an answer
// @Description Like an answer
// @Tags Answers
// @Accept json
// @Produce json
// @Security Bearer
// @Param answer_id path string true "Answer ID"
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /answer/{answer_id}/like [post]
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	var (
		dbLike  models.AnswerLike
		db      = database.GetDB()
		message string
	)

	answerID := r.PathValue("answer_id")
	userID := r.Context().Value("userID").(uuid.UUID)

	dbLike = models.AnswerLike{
		AnswerID: uuid.MustParse(answerID),
		UserID:   userID,
	}

	if err := db.Where("answer_id = ? AND user_id = ?", answerID, userID).First(&dbLike).Error; err != nil {
		if err := db.Create(&dbLike).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message = "Answer liked successfully"
	} else {
		if err := db.Delete(&dbLike).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message = "Answer unliked successfully"
	}

	response := models.NewBaseResponse(true, message, http.StatusOK, nil)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Edit Answer godoc
// @Summary Edit an answer
// @Description Edit an answer
// @Tags Answers
// @Accept json
// @Produce json
// @Security Bearer
// @Param answer_id path string true "Answer ID"
// @Param body body models.EditAnswerValidation true "Answer information"
// @Success 200 {object} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /answer/{answer_id} [put]
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	var (
		requestAnswer models.EditAnswerValidation
		dbAnswer      models.Answer
		db            = database.GetDB()
	)

	answerID := r.PathValue("answer_id")
	userID := r.Context().Value("userID").(uuid.UUID)

	if err := json.NewDecoder(r.Body).Decode(&requestAnswer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Where("id = ? AND user_id = ?", answerID, userID).First(&dbAnswer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbAnswer.EditAnswerCheckFields(requestAnswer)

	if err := db.Save(&dbAnswer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Answer edited successfully", http.StatusOK, dbAnswer.ToResponse())

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Delete Answer godoc
// @Summary Delete an answer
// @Description Delete an answer
// @Tags Answers
// @Accept json
// @Produce json
// @Security Bearer
// @Param answer_id path string true "Answer ID"
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /answer/{answer_id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		db       = database.GetDB()
		dbAnswer models.Answer
	)

	answerID := r.PathValue("answer_id")

	if err := db.Where("id = ?", answerID).First(&dbAnswer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Delete(&dbAnswer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Answer deleted successfully", http.StatusOK, nil)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}
