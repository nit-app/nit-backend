package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/requests"
)

func Register(engine *gin.Engine) {
	authGroup := engine.Group("/v1/auth")
	authGroup.POST("/sendCode", util.ValidateRequestData[requests.PhoneNumberRequest], SignIn)
	authGroup.POST("/confirm", util.ValidateRequestData[requests.OtpCheckRequest], CheckOTP)
	authGroup.GET("/revoke", Revoke)
}
