package events

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	_, err := env.DB().ExecContext(c,
		"insert into events (uuid, title, description, pricelow, pricehigh, agelimitlow, agelimithigh, location, ownerinfo, hasCertificate, plaindescription, isdraft, ismachinegenerated) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		header.UUID, header.Title, " ", header.PriceLow, header.PriceHigh, header.AgeLimitLow, header.AgeLimitHigh, header.Location, header.OwnerInfo, header.HasCertificate, header.PlainDescription, header.IsDraft, header.IsMachineGenerated)

	if err != nil {
		return nil, err
	}

	for _, tag := range header.Tags {
		appendTag := requests.AppendTag{EventUUID: header.UUID, Tag: tag}
		err := AppendTag(c, &appendTag)
		if err != nil {
			return nil, err
		}
	}

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
		join event_tags et on
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

	draftRow := env.DB().QueryRowContext(c, draftQuery, newUuid)

	draftHeader, err := ScanEventHeader(draftRow, nil)
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	return draftHeader, nil
}
