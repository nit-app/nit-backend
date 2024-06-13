package events

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/status"
)

func SetLinks(ctx context.Context, eventUUID string, links []*models.EventExternalLink) error {
	tx, err := env.DB().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	_, err = tx.ExecContext(ctx, "delete from event_external_links where \"eventUuid\" = $1", eventUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrappedErrors.New(status.InternalServerError, rollbackErr)
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	for _, link := range links {
		_, err := tx.ExecContext(ctx, "insert into event_external_links (\"linkUuid\", \"eventUuid\", title, url) values ($1, $2, $3, $4)",
			uuid.New(), eventUUID, link.Title, link.URL)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return wrappedErrors.New(status.InternalServerError, rollbackErr)
			}

			return wrapLinkInsertError(link.URL, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil
}

func wrapLinkInsertError(link string, err error) error {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return wrappedErrors.New(status.DuplicateValueEntry, errors.New("duplicate link: "+link)) // future check

		case "foreign_key_violation":
			return wrappedErrors.New(status.NoSuchEvent, err)
		}
	}

	return wrappedErrors.New(status.InternalServerError, err)
}
