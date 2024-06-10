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

	for _, day := range schedule {

		if err = checkSchedule(day); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return wrappedErrors.New(status.InternalServerError, rollbackErr)
			}

			return err
		}

		_, err := tx.ExecContext(ctx, "insert into event_schedule (scheduleUuid, \"eventUuid\", beginsAt, endsAt) values ($1, $2, $3, $4)",
			uuid.New(), eventUUID, day.BeginsAt, day.EndsAt)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return wrappedErrors.New(status.InternalServerError, rollbackErr)
			}

			return wrapScheduleInsertError(day, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil
}

func checkSchedule(schedule *models.EventSchedule) error {

	if schedule.BeginsAt.After(schedule.EndsAt) {

		return wrappedErrors.New(status.InvalidDataFormat, errors.New(fmt.Sprintf("event ends before it starts: %s - %s", schedule.BeginsAt, schedule.EndsAt)))
	}
	return nil
}

func wrapScheduleInsertError(schedule *models.EventSchedule, err error) error {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return wrappedErrors.New(status.TagAlreadySet, errors.New(fmt.Sprintf("duplicate schedule: %s - %s", schedule.BeginsAt, schedule.EndsAt))) // future check

		case "foreign_key_violation":
			return wrappedErrors.New(status.NoSuchEvent, err)
		}
	}

	return wrappedErrors.New(status.InternalServerError, err)
}
