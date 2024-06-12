package events

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services/events"
	"github.com/nit-app/nit-backend/services/fav"
)

func AddToFavorites(c *gin.Context) {
	req := c.Param("uuid")
	eventUUID, err := uuid.Parse(req)
	if err != nil {
		c.AbortWithStatusJSON(response.ErrorWithText(status.InvalidDataFormat, err.Error()))
		return
	}

	err = fav.AddToFavorites(c, util.GetUser(c).UUID, eventUUID.String())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(true))
}

func RemoveFromFavorites(c *gin.Context) {
	req := c.Param("uuid")
	eventUUID, err := uuid.Parse(req)
	if err != nil {
		c.AbortWithStatusJSON(response.ErrorWithText(status.InvalidDataFormat, err.Error()))
		return
	}

	err = fav.RemoveFromFavorites(c, util.GetUser(c).UUID, eventUUID.String())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(true))
}

func GetMyFavorites(c *gin.Context) {
	util.ServiceCall(c, events.GetUserFavorites, util.GetUser(c).UUID)
}
