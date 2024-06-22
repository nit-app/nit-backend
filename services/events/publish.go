package events

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/status"
	"go.uber.org/zap"
)

func Publish(ctx context.Context, uuid uuid.UUID) error {
	tx, err := env.DB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	event, err := txGetByUUID(ctx, tx, uuid)
	if err != nil {
		return err // already wrapped
	}

	if err := validateEvent(event); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "update events set isdraft = false, ismachinegenerated = false where uuid = $1", uuid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil
}

func validateEvent(event *models.Event) error {
	// todo add more checks

	if !event.IsDraft {
		return wrappedErrors.New(status.AlreadyPublished, errors.New("event already published"))
	}

	if len(event.Schedule) < 1 {
		return wrappedErrors.New(status.MalformedEvent, errors.New("empty schedule"))
	}

	if len(event.Tags) < 1 {
		return wrappedErrors.New(status.MalformedEvent, errors.New("empty tags"))
	}

	return nil
}
