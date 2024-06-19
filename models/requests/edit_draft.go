package requests

import "github.com/nit-app/nit-backend/models"

type EditDraft struct {
	EventUUID   string              `json:"eventUuid" binding:"required,uuid"`
	Object      *models.EventHeader `json:"object" binding:"required"`
	Description string              `json:"description"`
}
