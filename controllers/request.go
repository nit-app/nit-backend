package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
)

const requestDataKey = "request_data"

func ValidateRequestData[BodyType any]() func(c *gin.Context) {
	return func(c *gin.Context) {
		var body BodyType

		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatusJSON(response.ErrorWithText(status.InvalidDataFormat, err.Error()))
			return
		}

		c.Set(requestDataKey, body)
		c.Next()
	}
}

func GetRequestData[BodyType any](c *gin.Context) BodyType {
	return c.MustGet(requestDataKey).(BodyType)
}
