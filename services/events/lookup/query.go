package lookup

import (
	"context"
	"github.com/nit-app/nit-backend/env"
	"github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/services/events"
)

func Events(ctx context.Context, filters *requests.EventLookupFilters) ([]*responses.EventHeader, error) {
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
			es.beginsat,
			es.endsat,
			es.addedat,
			es.scheduleuuid,
			e.plainDescription
		from
			events e
		join event_tags et on
			e.uuid = et.uuid
		inner join event_schedule es on
			e.uuid = es."eventUuid"
			and es.beginsat >= $3
			and es.endsat <= $4
		where 
					e.deletedat is null
			and e.agelimitlow <= $1
			and e.pricelow <= $2
		group by
			e.uuid,
			e.favcount,
			es.scheduleuuid,
			es.addedat,
			es.beginsat,
			es.endsat
		order by
			e.favcount`

	ageLimitLowFilter := 100
	if filters.ExcludeAgeRestricted {
		ageLimitLowFilter = 0
	}

	priceFilter := 10000000
	if filters.ExcludePaid {
		priceFilter = 0
	}

	rows, err := env.DB().QueryContext(ctx, query, ageLimitLowFilter, priceFilter, filters.From, filters.To)
	if err != nil {
		return nil, errors.New(status.InternalServerError, err)
	}

	defer rows.Close()

	eventHeaders := make([]*responses.EventHeader, 0)
	for rows.Next() {
		eventObject, err := events.ScanEventHeader(rows, nil)
		if err != nil {
			return nil, errors.New(status.InternalServerError, err)
		}

		if !isMatchingTags(eventObject.Tags, filters.Tags) {
			continue
		}

		eventHeaders = append(eventHeaders, eventObject)
	}

	return eventHeaders, nil
}
