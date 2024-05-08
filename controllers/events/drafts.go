package events

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services/events"
)

func GetDrafts(c *gin.Context) {
	resp, err := events.GetDrafts(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(resp))
}

func CreateDraft(c *gin.Context) {
	req := util.GetRequestData[*models.EventHeader](c)
	header, err := events.CreateDraft(c, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(header))
}

func SetTags(c *gin.Context) {
	req := util.GetRequestData[*requests.SetTags](c)
	err := events.SetTags(c, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(true))
}
