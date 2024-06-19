package fav

import (
	"context"
	"database/sql"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/status"
	"go.uber.org/zap"
)

func RemoveFromFavorites(ctx context.Context, userUUID, eventUUID string) error {
	tx, err := env.DB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	res, err := tx.ExecContext(ctx, "delete from event_favorites where userUuid = $1 and eventUuid = $2", userUUID, eventUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	nrows, err := res.RowsAffected()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	if nrows > 0 {
		if err := addToCounterTx(ctx, tx, eventUUID, -1); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				zap.L().Error("cannot rollback", zap.Error(rollbackErr))
			}

			return wrapSetFavoriteError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil

}
