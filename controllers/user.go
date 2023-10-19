package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/sessions"
	"net/http"
)

type UserController struct {
}

func (uc *UserController) GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, response.Ok(sessions.Subject(c)))
}
