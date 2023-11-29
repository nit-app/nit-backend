package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/response"
)

func serviceCall[ArgType, RetType any](c *gin.Context, sm func(ctx context.Context, arg ArgType) (RetType, error), arg ArgType) {
	resp, err := sm(c, arg)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(response.Ok(resp))
}
