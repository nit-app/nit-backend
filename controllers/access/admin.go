package access

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/util"
)

func RequireAdmin(c *gin.Context) {
	if !util.IsAdmin(c) {
		c.AbortWithStatusJSON(response.ErrorWithText(errors.New(status.Forbidden, util.ErrNoAccess).MakeResponse()))
		return
	}

	c.Next()
}
