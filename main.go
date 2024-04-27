package main

import (
	"errors"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers"
	"github.com/nit-app/nit-backend/env"
	_ "github.com/nit-app/nit-backend/env/autoload"
	"github.com/nit-app/nit-backend/sessions"
	_ "github.com/nit-app/nit-backend/validators"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
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
	engine.Use(controllers.HandleErrors)
	engine.Use(controllers.CORS)

	controllers.Register(engine)

	server := &http.Server{
		Addr:    env.E().ListenAddress,
		Handler: engine,
	}

	sig := make(chan os.Signal, 1)
	go shutdown(sig, server)

	signal.Notify(sig, getShutdownSignals()...)

	zap.S().Infow("starting", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}

		zap.L().Fatal("listen error", zap.Error(err))
	}
}
