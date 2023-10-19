package sessions

import (
	"github.com/gin-gonic/gin"
)

func Current(c *gin.Context) *Session {
	session, ok := c.Get(SessionKey)
	if !ok {
		panic("session is supposed to be prefilled")
	}

	return session.(*Session)
}

func State(c *gin.Context) string {
	return Current(c).State
}

func Subject(c *gin.Context) *string {
	return Current(c).Subject
}

func SetAuthorized(session *Session, userUuid string) {
	session.State = StateAuthorized
	session.OTP = nil
	session.Subject = &userUuid

	session.Save()
}
