package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/services"
	"github.com/nit-app/nit-backend/sessions"
	"net/http"
)

type AuthController struct {
	AuthService *services.AuthService
}

type OtpDelegate interface {
	Start(session *sessions.Session, phoneNumber string) error
	CheckOTP(session *sessions.Session, otpCode string) error
	OtpCheckState() string
}

func (ac *AuthController) SignIn(c *gin.Context) {
	sendOtp(c, ac.AuthService)
}

func (ac *AuthController) CheckOTP(c *gin.Context) {
	checkOtp(c, ac.AuthService)
}

func (ac *AuthController) Revoke(c *gin.Context) {
	sessions.Current(c).Revoke()
	c.SetCookie(sessions.CookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, response.Ok(true))
}
