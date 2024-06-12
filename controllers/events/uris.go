package events

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/access"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/requests"
)

func Register(engine *gin.Engine) {
	eventsGroup := engine.Group("/v1/events")
	eventsGroup.POST("/lookup", util.ValidateRequestData[requests.EventLookupFilters], LookupEvents)
	eventsGroup.GET("/get/:uuid", GetEvent)

	eventsAuthorizedGroup := engine.Group("/v1/events")
	eventsAuthorizedGroup.Use(access.RequireAuth)
	eventsAuthorizedGroup.POST("/fav/:uuid/add", AddToFavorites)
	eventsAuthorizedGroup.POST("/fav/:uuid/remove", RemoveFromFavorites)
	eventsAuthorizedGroup.GET("/fav", GetMyFavorites)

	eventAdminGroup := engine.Group("/v1/eventAdmin")
	eventAdminGroup.Use(access.RequireAuth, access.RequireAdmin)

	eventAdminGroup.GET("/drafts", GetDrafts)
	eventAdminGroup.POST("/create", util.ValidateRequestData[*models.EventHeader], CreateDraft)
	eventAdminGroup.POST("/setTags", util.ValidateRequestData[*requests.SetTags], SetTags)
}
