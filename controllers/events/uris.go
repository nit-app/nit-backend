package events

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/requests"
)

func Register(engine *gin.Engine) {
	eventsGroup := engine.Group("/v1/events")
	eventsGroup.POST("/lookup", util.ValidateRequestData[requests.EventLookupFilters], LookupEvents)
	eventsGroup.GET("/get/:uuid", GetEvent)
}
