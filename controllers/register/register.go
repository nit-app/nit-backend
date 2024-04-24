package register

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/otp"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services/auth"
	"github.com/nit-app/nit-backend/sessions"
)

func Start(c *gin.Context) {
	otp.Send(c, auth.Register)
}

func CheckOTP(c *gin.Context) {
	otp.Check(c, auth.Register)
}

func Finish(c *gin.Context) {
	session := sessions.Current(c)

	if session.State != sessions.StateRegFinish {
		c.JSON(response.Error(status.BadFormState))
		return
	}

	req := util.GetRequestData[requests.FinishRegistrationRequest](c)

	_, err := auth.Register.Finish(session, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(response.ErrorWithText(status.BadRegistrationData, err.Error())) // catches db issues, reconsider logging
		return
	}

	c.JSON(response.Ok(true))
}
