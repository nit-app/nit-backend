package util

import (
	"context"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/services/user"
	"go.uber.org/zap"
)

const CtxKeyUser = "auth_user_uuid"

// GetUser extracts the current authenticated user from the request context and returns the user object.
// The user returned is always non-nil. GetUser panics if the user is not authenticated or subject cannot be queried
func GetUser(ctx context.Context) *models.User {
	val := ctx.Value(CtxKeyUser)
	if val == nil {
		panic("GetUser: no context key")
	}

	userUuid := val.(string)

	u, err := user.GetByUuid(ctx, userUuid)
	if err != nil {
		zap.L().Error("cannot receive user details", zap.Error(err))
		panic("GetUser: cannot receive user details from db")
	}

	return u
}
