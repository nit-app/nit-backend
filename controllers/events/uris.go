package events

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/access"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/requests"
)

func Register(engine *gin.Engine) {
	eventsGroup := engine.Group("/v1/events")
	eventsGroup.POST("/lookup", util.ValidateRequestData[requests.EventLookupFilters], LookupEvents)
	eventsGroup.GET("/get/:uuid", GetEvent)

	eventAdminGroup := engine.Group("/v1/eventAdmin")
	eventAdminGroup.Use(access.RequireAdmin)

	eventAdminGroup.GET("/drafts", GetDrafts)
	eventAdminGroup.POST("/create", CreateDraft)
	eventAdminGroup.POST("/appendTag", AppendTag)
}
