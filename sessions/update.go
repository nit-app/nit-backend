package sessions

import (
	"context"
	"github.com/nit-app/nit-backend/env"
	"time"
)

func (s *Session) Save() {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	params := []any{"state", s.State, "iat", s.Authorized}
	if s.Subject != nil {
		params = append(params, "sub", *s.Subject)
	}

	if s.OTP != nil {
		params = append(params, "phone", s.OTP.PhoneNumber, "code", s.OTP.Code, "attempt", s.OTP.Attempt)
	}

	if err := env.Redis().HSet(timeout, "session_"+s.tokHash, params).Err(); err != nil {
		panic(err)
	}
}
