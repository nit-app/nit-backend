package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/services"
	"github.com/nit-app/nit-backend/sessions"
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
