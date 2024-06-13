package events

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"strings"
)

func SetTags(c *gin.Context, req *requests.SetTags) error {
	tx, err := env.DB().BeginTx(c, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	_, err = tx.ExecContext(c, "delete from event_tags where uuid = $1", req.EventUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrappedErrors.New(status.InternalServerError, rollbackErr)
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	if err := setTagsTx(tx, c, req); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrappedErrors.New(status.InternalServerError, rollbackErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil
}

func setTagsTx(tx *sql.Tx, c *gin.Context, req *requests.SetTags) error {
	for _, tag := range req.Tags {
		tag = strings.ToLower(tag)

		_, err := tx.ExecContext(c, "insert into event_tags (uuid, tag) values ($1, $2)", req.EventUUID, tag)
		if err != nil {
			return wrapSetTagsError(tag, err)
		}
	}
	return nil
}

func wrapSetTagsError(tag string, err error) error {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return wrappedErrors.New(status.DuplicateValueEntry, errors.New("tag already exists: "+tag))
		case "foreign_key_violation":
			return wrappedErrors.New(status.NoSuchEvent, err)
		}
	}

	return wrappedErrors.New(status.InternalServerError, err)
}
