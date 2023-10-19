package main

import (
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers"
	_ "github.com/nit-app/nit-backend/env/autoload"
	"github.com/nit-app/nit-backend/services"
	"github.com/nit-app/nit-backend/services/otp"
	"github.com/nit-app/nit-backend/services/sms"
	"github.com/nit-app/nit-backend/sessions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	engine := gin.Default()
	engine.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(zap.L(), true))
	engine.Use(sessions.SessionKeeper)

	userService := &services.UserService{}
	otpService := &services.OtpService{Generator: otp.NewGenerator(), Carrier: sms.NewCarrier()}
	authService := &services.AuthService{OTP: otpService, UserService: userService}
	registerService := &services.RegisterService{OTP: otpService, UserService: userService}

	authController := &controllers.AuthController{AuthService: authService}
	userController := &controllers.UserController{}
	registerController := &controllers.RegisterController{RegisterService: registerService}

	authGroup := engine.Group("/v1/auth")
	authGroup.POST("/sendCode", authController.SignIn)
	authGroup.POST("/confirm", authController.CheckOTP)

	registerGroup := engine.Group("/v1/register")
	registerGroup.POST("/sendCode", registerController.StartRegistration)
	registerGroup.POST("/confirm", registerController.CheckOTP)
	registerGroup.POST("/finish", registerController.Finish)

	v1 := engine.Group("/v1")
	v1.Use(sessions.RequireAuth)

	v1.GET("/getMe", userController.GetMe)

	engine.Run(":8000")
}
