package models

import (
	"encoding/json"
	"net/http"
)

type (
	Response struct {
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		HttpStatus int         `json:"httpStatus"`
		Data       interface{} `json:"data,omitempty"`
		Errors     interface{} `json:"errors,omitempty"`
	}
)

func NewResponse(success bool, message string, httpStatus int, data interface{}, errors interface{}) Response {
	return Response{
		Success:    success,
		Message:    message,
		HttpStatus: httpStatus,
		Data:       data,
		Errors:     errors,
	}
}

func (b *Response) Write(w http.ResponseWriter) {
	w.WriteHeader(b.HttpStatus)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b.ToJson())
}

func (b Response) ToJson() []byte {
	jsonBytes, _ := json.Marshal(b)
	return jsonBytes
}
