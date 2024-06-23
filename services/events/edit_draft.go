package events

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
	"go.uber.org/zap"
)

func EditDraft(ctx context.Context, req *requests.EditDraft) (*models.EventHeader, error) {
	const query = `update events set (title, description, pricelow, pricehigh, 
		agelimitlow, agelimithigh, location, ownerinfo, hasCertificate, plaindescription,
		ismachinegenerated, modifiedat) = ($2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, current_timestamp) 
		where uuid = $1 and isdraft = true`

	tx, err := env.DB().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	header := req.Object

	res, err := tx.ExecContext(ctx, query,
		req.EventUUID, header.Title, req.Description, header.PriceLow, header.PriceHigh,
		header.AgeLimitLow, header.AgeLimitHigh, header.Location, header.OwnerInfo, header.HasCertificate,
		header.PlainDescription, header.IsMachineGenerated)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			return nil, wrappedErrors.New(status.InvalidDataFormatUnknown, err)
		}

		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	nrows, _ := res.RowsAffected()

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			zap.L().Error("cannot rollback", zap.Error(rollbackErr))
		}

		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	if nrows == 0 {
		return nil, wrappedErrors.New(status.NoSuchEvent, errors.New("no such event"))
	}

	return getDraftHeader(ctx, req.EventUUID)
}
