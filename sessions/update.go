package sessions

import (
	"context"
	"github.com/nit-app/nit-backend/env"
	"time"
)

const timeoutDuration = 3 * time.Second

func (s *Session) Save() {
	if len(s.tokHash) == 0 {
		panic("bad session object")
	}

	timeout, cancel := context.WithTimeout(context.Background(), timeoutDuration)
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
