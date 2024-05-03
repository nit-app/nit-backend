package util

import (
	"context"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/services/user"
	"go.uber.org/zap"
)

const CtxKeyUser = "auth_user_uuid"

// GetUser extracts the current authenticated user from the request context and returns the user object.
// GetUser returns nil if the user is not authenticated or subject cannot be queried
func GetUser(ctx context.Context) *models.User {
	val := ctx.Value(CtxKeyUser)
	if val == nil {
		return nil
	}

	userUuid := val.(string)

	u, err := user.GetByUuid(ctx, userUuid)
	if err != nil {
		zap.L().Error("cannot receive user details", zap.Error(err))
		return nil
	}

	return u
}
