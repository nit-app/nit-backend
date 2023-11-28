package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services"
	"github.com/nit-app/nit-backend/sessions"
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
		c.JSON(response.Error(status.BadFormState))
		return
	}

	req := GetRequestData[requests.FinishRegistrationRequest](c)

	_, err := rc.RegisterService.Finish(session, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(response.ErrorWithText(status.BadRegistrationData, err.Error())) // catches db issues, reconsider logging
		return
	}

	c.JSON(response.Ok(true))
}
