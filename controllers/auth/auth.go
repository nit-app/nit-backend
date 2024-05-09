package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/otp"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services/auth"
	"github.com/nit-app/nit-backend/sessions"
)

func SignIn(c *gin.Context) {
	otp.Send(c, auth.SignIn)
}

func CheckOTP(c *gin.Context) {
	otp.Check(c, auth.SignIn)
}

func Revoke(c *gin.Context) {
	sessions.Current(c).Revoke()
	c.SetCookie(sessions.CookieName, "", -1, "/", "", false, true)
	c.JSON(response.Ok(true))
}
