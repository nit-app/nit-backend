package response

import (
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/nit-app/nit-backend/models/status"
	"net/http"
	"time"
)

const defaultStatus = http.StatusBadRequest

func Ok[T any](object T) *responses.BaseResponse[T] {
	return &responses.BaseResponse[T]{Object: object, Timestamp: time.Now(), Status: http.StatusOK}
}

func Error(code string) (int, *responses.ErrorResponse) {
	return ErrorWithText(code, "")
}

func ErrorWithText(code, text string) (int, *responses.ErrorResponse) {
	resp := &responses.ErrorResponse{Code: code}
	statusCode := getStatus(code)
	resp.BaseResponse = responses.BaseResponse[*struct{}]{Timestamp: time.Now(), Status: statusCode, Text: text}
	return statusCode, resp
}

func getStatus(code string) int {
	val, ok := status.Codes[code]
	if !ok {
		return defaultStatus
	}
	return val
}
