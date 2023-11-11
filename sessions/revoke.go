package sessions

import (
	"context"
	"github.com/nit-app/nit-backend/env"
)

// Revoke deletes current session, drops the object and invalidates the session token.
// The session object is not usable after revocation and cannot be restored.
// Do not call to Save after revoking.
func (s *Session) Revoke() {
	if len(s.tokHash) == 0 {
		panic("session object is already revoked or invalid")
	}

	timeout, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	env.Redis().Del(timeout, "session_"+s.tokHash)
}
