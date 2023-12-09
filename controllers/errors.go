package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"go.uber.org/zap"
)

func HandleErrors(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	if len(c.Errors) > 1 {
		zap.L().Error("error handling middleware cannot handle more than one error at once")
		return
	}

	var serviceError *wrappedErrors.Error
	ok := errors.As(c.Errors.Last(), &serviceError)
	if !ok {
		zap.L().Warn("unsolicited non-wrapped error", zap.Error(c.Errors[0]))
		return
	}

	if serviceError.Type() == status.InternalServerError {
		zap.L().Error("internal service error", zap.Error(serviceError.Unwrap()))
	}

	c.AbortWithStatusJSON(response.ErrorWithText(serviceError.MakeResponse()))
}
