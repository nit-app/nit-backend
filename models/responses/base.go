package responses

import "time"

type BaseResponse[T any] struct {
	Timestamp time.Time `json:"timestamp"`
	Object    T         `json:"object,omitempty"`
	Status    int       `json:"status"`
	Text      string    `json:"text,omitempty"`
}
