package events

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/status"
)

func CreateDraft(c *gin.Context, header *models.EventHeader) (*models.EventHeader, error) {
	header.IsDraft = true

	newUuid := uuid.New()
	header.UUID = newUuid.String()

	tx, err := env.DB().BeginTx(c, &sql.TxOptions{})
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	_, err = tx.ExecContext(c, "insert into events (uuid, title, description, pricelow, pricehigh, agelimitlow, agelimithigh, location, ownerinfo, hasCertificate, plaindescription, isdraft, ismachinegenerated) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		header.UUID, header.Title, " ", header.PriceLow, header.PriceHigh, header.AgeLimitLow, header.AgeLimitHigh, header.Location, header.OwnerInfo, header.HasCertificate, header.PlainDescription, header.IsDraft, header.IsMachineGenerated)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, wrappedErrors.New(status.InternalServerError, err)
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			return nil, wrappedErrors.New(status.InvalidDataFormat, err)
		}

		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	setTags := &requests.SetTags{EventUUID: header.UUID, Tags: header.Tags}
	if err := setTagsTx(tx, c, setTags); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, wrappedErrors.New(status.InternalServerError, err)
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	return getDraftHeader(c, newUuid.String())
}

func getDraftHeader(c *gin.Context, uuid string) (*models.EventHeader, error) {
	const draftQuery = `
		select
			e.uuid,
			e.title,
			e.priceLow,
			e.priceHigh,
			e.ageLimitLow,
			e.ageLimitHigh,
			e.location,
			e.ownerInfo,
			string_agg(distinct et.tag, ',') as tags,
			e.createdat,
			e.modifiedat,
			es.beginsat,
			es.endsat,
			es.addedat,
			es.scheduleuuid,
			e.hascertificate,
			e.plainDescription,
			e.favcount,
			e.isDraft,
			e.isMachineGenerated
		from
			events e
		left join event_tags et on
			e.uuid = et.uuid
		left join event_schedule es on
			e.uuid = es."eventUuid"
		where
		    e.uuid = $1
			and e.deletedat is null
		group by
			e.uuid,
			e.favcount,
			es.scheduleuuid,
			es.addedat,
			es.beginsat,
			es.endsat`

	draftRow := env.DB().QueryRowContext(c, draftQuery, uuid)

	draftHeader, err := ScanEventHeader(draftRow, nil)
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	return draftHeader, nil
}
