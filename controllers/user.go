package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/response"
)

func GetMe(c *gin.Context) {
	c.JSON(response.Ok(util.GetUser(c)))
}
