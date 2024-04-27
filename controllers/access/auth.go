package access

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/controllers/util"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/response"
	"github.com/nit-app/nit-backend/sessions"
)

func RequireAuth(c *gin.Context) {
	sessionRaw, ok := c.Get(sessions.SessionKey)
	if !ok || sessionRaw.(*sessions.Session).State != sessions.StateAuthorized || sessionRaw.(*sessions.Session).Subject == nil {
		c.AbortWithStatusJSON(response.Error(status.Unauthorized))
		return
	}

	c.Set(util.CtxKeyUser, *sessionRaw.(*sessions.Session).Subject)
	c.Next()
}
