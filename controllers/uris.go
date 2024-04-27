package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/access"
	"github.com/nit-app/nit-backend/controllers/auth"
	"github.com/nit-app/nit-backend/controllers/events"
	"github.com/nit-app/nit-backend/controllers/register"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(engine *gin.Engine) {
	engine.StaticFile("/docs.yaml", "schema/docs.yaml")
	engine.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../docs.yaml")))

	auth.Register(engine)
	register.Register(engine)
	events.Register(engine)

	v1 := engine.Group("/v1")
	v1.Use(access.RequireAuth)

	v1.GET("/getMe", GetMe)
}
