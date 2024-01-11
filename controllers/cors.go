package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllowedOrigin(c *gin.Context) string {
	origin := c.GetHeader("Origin")
	if len(origin) == 0 {
		origin = c.GetHeader("Referer")
	}
	if len(origin) == 0 {
		origin = "*"
	}

	return origin
}

func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", getAllowedOrigin(c))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
