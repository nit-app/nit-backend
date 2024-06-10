package requests

import "github.com/nit-app/nit-backend/models"

type SetSchedule struct {
	EventUUID string                  `json:"eventUuid" binding:"required,uuid"`
	Schedule  []*models.EventSchedule `json:"schedule" binding:"required,dive"`
}
