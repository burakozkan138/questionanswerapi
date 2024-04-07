package models

import "encoding/json"

type (
	BaseResponse struct {
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		HttpStatus int         `json:"httpStatus"`
		Data       interface{} `json:"data,omitempty"`
	}

	ErrorResponse struct {
		BaseResponse
		Errors map[string]string `json:"errors"`
	}
)

func NewBaseResponse(success bool, message string, httpStatus int, data interface{}) BaseResponse {
	return BaseResponse{
		Success:    success,
		Message:    message,
		HttpStatus: httpStatus,
		Data:       data,
	}
}

func NewErrorResponse(success bool, message string, httpStatus int, data interface{}, errors map[string]string) ErrorResponse {
	return ErrorResponse{
		BaseResponse: BaseResponse{
			Success:    success,
			Message:    message,
			HttpStatus: httpStatus,
			Data:       data,
		},
		Errors: errors,
	}
}

func (b BaseResponse) ToJson() []byte {
	jsonBytes, _ := json.Marshal(b)
	return jsonBytes
}

func (e ErrorResponse) Error() []byte {
	jsonBytes, _ := json.Marshal(e)
	return jsonBytes
}
