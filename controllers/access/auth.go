package access

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/sessions"
)

func RequireAuth(c *gin.Context) {
	_, ok := c.Get(util.CtxKeyUser)
	if !ok {
		c.AbortWithStatusJSON(response.Error(status.Unauthorized))
		return
	}
	c.Next()
}

func CheckAuth(c *gin.Context) {
	defer c.Next()
	sessionRaw, ok := c.Get(sessions.SessionKey)
	if !ok || sessionRaw.(*sessions.Session).State != sessions.StateAuthorized || sessionRaw.(*sessions.Session).Subject == nil {
		return
	}

	c.Set(util.CtxKeyUser, *sessionRaw.(*sessions.Session).Subject)
}
