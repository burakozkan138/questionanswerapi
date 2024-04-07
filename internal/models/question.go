package models

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Question struct {
		BaseModel
		Title   string         `json:"title"`
		Content string         `json:"content"`
		Slug    string         `json:"slug"`
		UserID  uuid.UUID      `json:"user_id"`
		User    User           `json:"user"`
		Answers []Answer       `json:"answers,omitempty"`
		Likes   []QuestionLike `json:"likes,omitempty"`
	}

	QuestionResponse struct {
		BaseModel
		Title        string           `json:"title"`
		Content      string           `json:"content"`
		Slug         string           `json:"slug"`
		User         UserInform       `json:"user"`
		Answers      []AnswerResponse `json:"answers,omitempty"`
		AnswersCount int              `json:"answers_count"`
		Likes        int              `json:"likes_count"`
	}

	CreateQuestionValidation struct {
		Title   string `json:"title" validate:"required,min=5,max=100"`
		Content string `json:"content" validate:"required,min=5,max=500"`
	}

	EditQuestionValidation struct {
		Title   string `json:"title,omitempty" validate:"omitempty,min=5,max=100"`
		Content string `json:"content,omitempty" validate:"omitempty,min=5,max=500"`
	}

	QuestionLike struct {
		BaseModel
		UserID     uuid.UUID `json:"user_id"`
		QuestionID uuid.UUID `json:"question_id"`
	}
)

func (q *Question) BeforeSave(tx *gorm.DB) (err error) {
	if len(q.Title) > 0 {
		q.Slug = q.GenerateSluq()
	}

	return
}

func (q *Question) GenerateSluq() string {
	return strings.ReplaceAll(strings.ToLower(q.Title), " ", "-")
}

func (q *Question) TableName() string {
	return "questions"
}

func (q *Question) ToResponse() QuestionResponse {
	Answers := make([]AnswerResponse, 0)
	for _, a := range q.Answers {
		if a.User.ID == uuid.Nil { // sadece answers sayısı lazım :>
			continue
		}
		Answers = append(Answers, a.ToResponse())
	}

	response := QuestionResponse{
		BaseModel:    q.BaseModel,
		Title:        q.Title,
		Content:      q.Content,
		Slug:         q.Slug,
		User:         q.User.ToUserInformResponse(),
		Answers:      Answers,
		Likes:        len(q.Likes),
		AnswersCount: len(q.Answers),
	}

	return response
}

func MapToQuestionResponse(questions []Question) []QuestionResponse {
	response := make([]QuestionResponse, len(questions))

	for i, question := range questions {
		response[i] = question.ToResponse()
	}

	return response
}

func (q *Question) EditQuestionCheckFields(EditQuestion EditQuestionValidation) {
	if len(EditQuestion.Title) > 0 {
		q.Title = EditQuestion.Title
	}

	if len(EditQuestion.Content) > 0 {
		q.Content = EditQuestion.Content
	}
}
