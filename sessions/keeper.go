package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/nit-app/nit-backend/response"
	"net/http"
)

const (
	CookieName = "isomiso"
	ttlSecs    = 60 * 60 * 24 * 14
)

func SessionKeeper(c *gin.Context) {
	defer c.Next()

	if _, ok := c.Get(SessionKey); ok {
		return
	}

	token, err := c.Cookie(CookieName)

	if err != nil {
		token = resetSession(c)
	}

	s := getSessionOrReset(c, token)
	c.Set(SessionKey, s)
}

func RequireAuth(c *gin.Context) {
	sessionRaw, ok := c.Get(SessionKey)
	if !ok || sessionRaw.(*Session).State != StateAuthorized {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	c.Next()
}

func resetSession(c *gin.Context) string {
	token, err := createSession()
	if err != nil {
		panic(err)
	}

	c.SetCookie(CookieName, token, ttlSecs, "/", "", false, true)
	return token
}

func getSessionOrReset(c *gin.Context, token string) *Session {
	s, err := getSession(token)
	if err == nil {
		return s
	}

	if err == ErrNoSuchSession {
		token = resetSession(c)

		s, err := getSession(token)
		if err != nil {
			panic(err)
		}

		return s
	}

	panic(err)
}
