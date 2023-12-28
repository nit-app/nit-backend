package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/services/events/get_event"
	"github.com/nit-app/nit-backend/services/events/lookup"
)

func LookupEvents(c *gin.Context) {
	req := GetRequestData[requests.EventLookupFilters](c)
	serviceCall(c, lookup.Events, &req)
}

func GetEvent(c *gin.Context) {
	req := GetRequestData[requests.EventIdRequest](c)
	serviceCall(c, get_event.Event, &req)
}
