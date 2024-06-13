package events

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/status"
)

func SetSchedule(ctx context.Context, eventUUID string, schedule []*models.EventSchedule) error {
	tx, err := env.DB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	_, err = tx.ExecContext(ctx, "delete from event_schedule where \"eventUuid\" = $1", eventUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrappedErrors.New(status.InternalServerError, rollbackErr)
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	for _, item := range schedule {
		if err := insertScheduleItem(ctx, tx, eventUUID, item); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return wrappedErrors.New(status.InternalServerError, rollbackErr)
			}

			return wrapScheduleInsertError(item, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil
}

func insertScheduleItem(ctx context.Context, tx *sql.Tx, eventUUID string, item *models.EventSchedule) error {
	var err error

	if err = checkSchedule(item); err == nil {
		_, err = tx.ExecContext(ctx, "insert into event_schedule (scheduleUuid, \"eventUuid\", beginsAt, endsAt) values ($1, $2, $3, $4)",
			uuid.New(), eventUUID, item.BeginsAt, item.EndsAt)
	}

	return err
}

func checkSchedule(item *models.EventSchedule) error {
	if item.BeginsAt.After(item.EndsAt) {
		return wrappedErrors.New(status.InvalidDataFormat, fmt.Errorf("event ends before it starts: %s - %s", item.BeginsAt, item.EndsAt))
	}

	return nil
}

func wrapScheduleInsertError(item *models.EventSchedule, err error) error {
	var pqErr *pq.Error
	var formatErr *wrappedErrors.Error

	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return wrappedErrors.New(status.DuplicateValueEntry, fmt.Errorf("duplicate schedule: %s - %s", item.BeginsAt, item.EndsAt))
		case "foreign_key_violation":
			return wrappedErrors.New(status.NoSuchEvent, err)
		}
	} else if errors.As(err, &formatErr) {
		return err
	}

	return wrappedErrors.New(status.InternalServerError, err)
}
