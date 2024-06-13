package events

import (
	"context"
	"github.com/nit-app/nit-backend/env"
	"github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/status"
)

func GetDrafts(ctx context.Context) ([]*models.EventHeader, error) {
	const query = `
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
			coalesce(es.beginsat, '1970-01-01'),
			coalesce(es.endsat, '1970-01-01'),
			coalesce(es.addedat, '1970-01-01'),
			coalesce(es.scheduleuuid, '00000000-0000-0000-0000-000000000000'),
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
		    e.isDraft
			and e.deletedat is null
		group by
			e.uuid,
			e.favcount,
			es.scheduleuuid,
			es.addedat,
			e.modifiedat,
			es.beginsat,
			es.endsat
		order by e.modifiedat desc`

	rows, err := env.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New(status.InternalServerError, err)
	}

	defer rows.Close()

	eventHeaders := make([]*models.EventHeader, 0)
	for rows.Next() {
		eventObject, err := ScanEventHeader(rows, nil)
		if err != nil {
			return nil, errors.New(status.InternalServerError, err)
		}

		eventHeaders = append(eventHeaders, eventObject)
	}

	return eventHeaders, nil
}
