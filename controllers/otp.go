package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/sessions"
	"net/http"
)

func sendOtp(c *gin.Context, delegate OtpDelegate) {
	if sessions.State(c) != sessions.StateUnauthorized {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "bad sign in state"))
		return
	}

	var req requests.PhoneNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	if err := delegate.Start(sessions.Current(c), req.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Ok(true))
}

func checkOtp(c *gin.Context, delegate OtpDelegate) {
	if sessions.State(c) != delegate.OtpCheckState() {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "bad sign in state"))
		return
	}

	var req requests.OtpCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	if err := delegate.CheckOTP(sessions.Current(c), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Ok(true))
}
