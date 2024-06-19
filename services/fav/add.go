package fav

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/status"
	"go.uber.org/zap"
)

func AddToFavorites(ctx context.Context, userUUID, eventUUID string) error {
	tx, err := env.DB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	_, err = tx.ExecContext(ctx, "insert into event_favorites (userUuid, eventUuid) values ($1, $2)", userUUID, eventUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		return wrapSetFavoriteError(err)
	}

	if err := addToCounterTx(ctx, tx, eventUUID, 1); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil
}

func wrapSetFavoriteError(err error) error {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return nil // not an error to like an event twice

		case "foreign_key_violation":
			return wrappedErrors.New(status.NoSuchEvent, err)
		}
	}

	return wrappedErrors.New(status.InternalServerError, err)
}

func addToCounterTx(ctx context.Context, tx *sql.Tx, eventUUID string, delta int) error {
	_, err := tx.ExecContext(ctx, "update events set favcount = favcount + $2 where \"uuid\" = $1", eventUUID, delta)
	return err
}
