package register

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/requests"
)

func Register(engine *gin.Engine) { // this method name actually isn't related to package's purpose
	registerGroup := engine.Group("/v1/register")
	registerGroup.POST("/sendCode", util.ValidateRequestData[requests.PhoneNumberRequest], Start)
	registerGroup.POST("/confirm", util.ValidateRequestData[requests.OtpCheckRequest], CheckOTP)
	registerGroup.POST("/finish", util.ValidateRequestData[requests.FinishRegistrationRequest], Finish)
}
