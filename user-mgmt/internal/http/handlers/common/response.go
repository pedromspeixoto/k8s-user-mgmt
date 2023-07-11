package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	StatusCode   int    `json:"status_code,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", r.StatusCode, r.ErrorMessage)
}

func Json(w http.ResponseWriter, httpCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	res := Response{
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(res)
}

func Text(w http.ResponseWriter, httpCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpCode)
	w.Write([]byte(message))
}

func Err(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	res := ErrorResponse{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
	json.NewEncoder(w).Encode(res)
}
