package requests

import "github.com/nit-app/nit-backend/models"

type SetLinks struct {
	EventUUID string                      `json:"eventUuid" binding:"required,uuid"`
	Links     []*models.EventExternalLink `json:"links" binding:"required,dive"`
}
