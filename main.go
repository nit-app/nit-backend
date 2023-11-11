package main

import (
	"errors"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers"
	"github.com/nit-app/nit-backend/env"
	_ "github.com/nit-app/nit-backend/env/autoload"
	"github.com/nit-app/nit-backend/services"
	"github.com/nit-app/nit-backend/services/otp"
	"github.com/nit-app/nit-backend/services/sms"
	"github.com/nit-app/nit-backend/sessions"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer env.Shutdown()

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	engine := gin.Default()
	engine.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(zap.L(), true))

	engine.Use(sessions.SessionKeeper)
	setCors(engine)

	engine.StaticFile("/docs.yaml", "schema/docs.yaml")
	engine.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../docs.yaml")))

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
	authGroup.GET("/revoke", authController.Revoke)

	registerGroup := engine.Group("/v1/register")
	registerGroup.POST("/sendCode", registerController.StartRegistration)
	registerGroup.POST("/confirm", registerController.CheckOTP)
	registerGroup.POST("/finish", registerController.Finish)

	v1 := engine.Group("/v1")
	v1.Use(sessions.RequireAuth)

	v1.GET("/getMe", userController.GetMe)

	server := &http.Server{
		Addr:    env.E().ListenAddress,
		Handler: engine,
	}

	sig := make(chan os.Signal, 1)
	go shutdown(sig, server)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGKILL)

	zap.S().Infow("starting", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}

		zap.L().Fatal("listen error", zap.Error(err))
	}
}
