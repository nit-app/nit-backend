package util

import (
	"context"
	"errors"
	"github.com/nit-app/nit-backend/controllers/util"
)

var ErrNoAccess = errors.New("forbidden")

func IsAdmin(ctx context.Context) bool {
	u := util.GetUser(ctx)

	if u == nil { // unauthorized
		return false
	}

	return u.IsAdmin
}
