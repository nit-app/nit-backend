package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/sessions"
)

type Delegate interface {
	Start(session *sessions.Session, phoneNumber string) error
	CheckOTP(session *sessions.Session, otpCode string) error
}

func Send(c *gin.Context, delegate Delegate) {
	if sessions.State(c) != sessions.StateUnauthorized {
		c.JSON(response.Error(status.BadFormState))
		return
	}

	req := util.GetRequestData[requests.PhoneNumberRequest](c)

	if err := delegate.Start(sessions.Current(c), req.PhoneNumber); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(true))
}

func Check(c *gin.Context, delegate Delegate) {
	req := util.GetRequestData[requests.OtpCheckRequest](c)

	if err := delegate.CheckOTP(sessions.Current(c), req.Code); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(true))
}
