package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/services/events/lookup"
)

func LookupEvents(c *gin.Context) {
	req := GetRequestData[requests.EventLookupFilters](c)
	serviceCall(c, lookup.Events, &req)
}
