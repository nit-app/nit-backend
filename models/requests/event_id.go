package requests

import "github.com/google/uuid"

type EventIdRequest struct {
	UUID uuid.UUID `json:"uuid" binding:"required"`
}
