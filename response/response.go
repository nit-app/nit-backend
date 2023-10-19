package response

import (
	"github.com/nit-app/nit-backend/models/responses"
	"net/http"
	"time"
)

func Ok[T any](object T) *responses.BaseResponse[T] {
	return &responses.BaseResponse[T]{Object: object, Timestamp: time.Now(), Status: http.StatusOK}
}

func Error(status int, text string) *responses.BaseResponse[*struct{}] {
	return &responses.BaseResponse[*struct{}]{Timestamp: time.Now(), Status: status, Text: text}
}
