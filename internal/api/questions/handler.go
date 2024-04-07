package questions

import (
	"burakozkan138/questionanswerapi/internal/database"
	"burakozkan138/questionanswerapi/internal/models"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct{}

// Fetch all questions godoc
// @Summary Fetch all questions
// @Description Fetch all questions
// @Tags Questions
// @Accept json
// @Produce json
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /question/ [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	var (
		dbQuestions []models.Question
		db          = database.GetDB()
	)

	if err := db.Preload("Likes").Preload("User").Preload("Answers").Find(&dbQuestions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Questions fetched successfully", http.StatusOK, models.MapToQuestionResponse(dbQuestions))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Fetch question by ID godoc
// @Summary Fetch question by ID
// @Description Fetch question by ID
// @Tags Questions
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /question/{question_id} [get]
func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	var (
		dbQuestion models.Question
		db         = database.GetDB()
	)

	questionID := r.PathValue("question_id")

	if err := db.Preload("User").Preload("Likes").Preload("Answers.User").Preload("Answers.Likes").Where("id = ?", questionID).First(&dbQuestion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Question fetched successfully", http.StatusOK, dbQuestion.ToResponse())

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Edit question godoc
// @Summary Edit question
// @Description Edit question
// @Tags Questions
// @Accept json
// @Produce json
// @Security Bearer
// @Param question_id path string true "Question ID"
// @Param body body models.EditQuestionValidation true "Question information"
// @Success 200 {object} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /question/{question_id}/edit [put]
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	var (
		requestQuestion models.EditQuestionValidation
		dbQuestion      models.Question
		db              = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestQuestion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	questionID := r.PathValue("question_id")
	userID := r.Context().Value("userID").(uuid.UUID)

	if err := db.Where("id = ? AND user_id = ?", questionID, userID).Preload("User").First(&dbQuestion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbQuestion.EditQuestionCheckFields(requestQuestion)

	if err := db.Save(&dbQuestion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Question edit successfully", http.StatusOK, dbQuestion.ToResponse())

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Like question godoc
// @Summary Like question
// @Description Like question
// @Tags Questions
// @Accept json
// @Produce json
// @Security Bearer
// @Param question_id path string true "Question ID"
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /question/{question_id}/like [post]
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	var (
		dbLike  models.QuestionLike
		db      = database.GetDB()
		message string
	)

	questionID := r.PathValue("question_id")
	userID := r.Context().Value("userID").(uuid.UUID)

	dbLike = models.QuestionLike{
		QuestionID: uuid.MustParse(questionID),
		UserID:     userID,
	}

	if err := db.Where("question_id = ? AND user_id = ?", questionID, userID).First(&dbLike).Error; err != nil {

		if err := db.Create(&dbLike).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message = "Question liked successfully"
	} else {
		if err := db.Delete(&dbLike).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message = "Question unliked successfully"
	}

	response := models.NewBaseResponse(true, message, http.StatusOK, nil)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}

// Ask godoc
// @Summary Ask a question
// @Description Ask a question
// @Tags Questions
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body models.CreateQuestionValidation true "Question information"
// @Success 201 {object} models.BaseResponse
// @Failure 400 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /question/ask [post]
func (h *Handler) Ask(w http.ResponseWriter, r *http.Request) {
	var (
		requestQuestion models.CreateQuestionValidation
		dbQuestion      models.Question
		db              = database.GetDB()
	)

	if err := json.NewDecoder(r.Body).Decode(&requestQuestion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")

	dbQuestion = models.Question{
		Title:   requestQuestion.Title,
		Content: requestQuestion.Content,
		UserID:  userID.(uuid.UUID),
	}

	if err := db.Create(&dbQuestion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Preload("User").First(&dbQuestion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Question created successfully", http.StatusCreated, dbQuestion.ToResponse())

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusCreated)
}

// Delete godoc
// @Summary Delete question
// @Description Delete question
// @Tags Questions
// @Accept json
// @Produce json
// @Security Bearer
// @Param question_id path string true "Question ID"
// @Success 200 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /question/{question_id}/delete [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		dbQuestion models.Question
		db         = database.GetDB()
	)

	questionID := r.PathValue("question_id")
	userID := r.Context().Value("userID").(uuid.UUID)

	if err := db.Where("id = ? AND user_id = ?", questionID, userID).First(&dbQuestion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Delete(&dbQuestion).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewBaseResponse(true, "Question deleted successfully", http.StatusOK, nil)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	w.WriteHeader(http.StatusOK)
}
