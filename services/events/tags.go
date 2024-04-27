package events

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"strings"
)

func AppendTag(c *gin.Context, req *requests.AppendTag) error {
	req.Tag = strings.ToLower(req.Tag)

	_, err := env.DB().ExecContext(c, "insert into event_tags (uuid, tag) values ($1, $2)", req.EventUUID, req.Tag)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return wrappedErrors.New(status.TagAlreadySet, errors.New("tag already exists: "+req.Tag))
			case "foreign_key_violation":
				return wrappedErrors.New(status.NoSuchEvent, err)
			}
		}

		return wrappedErrors.New(status.InternalServerError, err)
	}

	return nil
}
