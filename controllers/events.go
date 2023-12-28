package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services/events"
	"github.com/nit-app/nit-backend/services/events/lookup"
)

func LookupEvents(c *gin.Context) {
	req := GetRequestData[requests.EventLookupFilters](c)
	serviceCall(c, lookup.Events, &req)
}

func GetEvent(c *gin.Context) {
	req := c.Param("uuid")
	eventUUID, err := uuid.Parse(req)
	if err != nil {
		c.AbortWithStatusJSON(response.ErrorWithText(status.InvalidDataFormat, err.Error()))
		return
	}
	serviceCall(c, events.GetByUUID, eventUUID)
}
