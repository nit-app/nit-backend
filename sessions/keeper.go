package sessions

import (
	"errors"
	"github.com/gin-gonic/gin"
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

	if errors.Is(err, ErrNoSuchSession) {
		token = resetSession(c)

		s, err := getSession(token)
		if err != nil {
			panic(err)
		}

		return s
	}

	panic(err)
}
