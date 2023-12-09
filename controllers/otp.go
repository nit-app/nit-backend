package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/sessions"
)

func sendOtp(c *gin.Context, delegate OtpDelegate) {
	if sessions.State(c) != sessions.StateUnauthorized {
		c.JSON(response.Error(status.BadFormState))
		return
	}

	req := GetRequestData[requests.PhoneNumberRequest](c)

	if err := delegate.Start(sessions.Current(c), req.PhoneNumber); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(true))
}

func checkOtp(c *gin.Context, delegate OtpDelegate) {
	req := GetRequestData[requests.OtpCheckRequest](c)

	if err := delegate.CheckOTP(sessions.Current(c), req.Code); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(true))
}
