package models

import (
	"github.com/google/uuid"
)

type (
	Answer struct {
		BaseModel
		Content string `json:"content"`

		QuestionID uuid.UUID    `json:"question_id"`
		Question   Question     `json:"question"`
		UserID     uuid.UUID    `json:"user_id"`
		User       User         `json:"user"`
		Likes      []AnswerLike `json:"likes"`
	}

	AnswerLike struct {
		BaseModel
		UserID   uuid.UUID `json:"user_id"`
		AnswerID uuid.UUID `json:"answer_id"`
	}

	CreateAnswerValidation struct {
		Content string `json:"content" validate:"required,min=10,max=500"`
	}

	EditAnswerValidation struct {
		Content string `json:"content" validate:"omitempty,min=10,max=500"`
	}

	AnswerResponse struct {
		BaseModel
		Content string     `json:"content"`
		User    UserInform `json:"user"`
		Likes   int        `json:"likes"`
	}
)

func (a *Answer) ToResponse() AnswerResponse {
	return AnswerResponse{
		BaseModel: a.BaseModel,
		Content:   a.Content,
		User:      a.User.ToUserInformResponse(),
		Likes:     len(a.Likes),
	}
}

func MapToAnswerResponse(answers []Answer) []AnswerResponse {
	response := make([]AnswerResponse, len(answers))

	for i, answer := range answers {
		response[i] = answer.ToResponse()
	}

	return response
}

func (a *Answer) EditAnswerCheckFields(request EditAnswerValidation) {
	if request.Content != "" {
		a.Content = request.Content
	}
}
