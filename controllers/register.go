package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services"
	"github.com/nit-app/nit-backend/sessions"
	"net/http"
)

type RegisterController struct {
	RegisterService *services.RegisterService
}

func (rc *RegisterController) StartRegistration(c *gin.Context) {
	sendOtp(c, rc.RegisterService)
}

func (rc *RegisterController) CheckOTP(c *gin.Context) {
	checkOtp(c, rc.RegisterService)
}

func (rc *RegisterController) Finish(c *gin.Context) {
	session := sessions.Current(c)

	if session.State != sessions.StateRegFinish {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "bad register state"))
		return
	}

	var req requests.FinishRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	userUuid, err := rc.RegisterService.Finish(session, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error())) // catches db issues, reconsider logging
		return
	}

	sessions.SetAuthorized(session, userUuid)

	c.JSON(http.StatusOK, response.Ok(true))
}
