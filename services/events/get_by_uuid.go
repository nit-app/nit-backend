package events

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/nit-app/nit-backend/models/status"
)

func GetByUUID(ctx context.Context, uuid uuid.UUID) (*responses.Event, error) {
	const headerQuery = `
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
			e.plainDescription,
			e.description
		from
			events e
		join event_tags et on
			e.uuid = et.uuid
		inner join event_schedule es on
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

	const scheduleQuery = `
		select
			es.beginsat,
			es.endsat,
			es.addedat,
			es.scheduleuuid
		from
			event_schedule es
		where 
		    es."eventUuid" = $1`

	const linksQuery = `
		select
    		el."linkUuid",
    		el.title,
    		el.url,
    		el.addedAt
		from 
		    event_external_links el
		where 
		    el."eventUuid" = $1`

	headerRow := env.DB().QueryRowContext(ctx, headerQuery, uuid)

	event := &responses.Event{}

	eventHeader, err := ScanEventHeader(headerRow, &event.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, wrappedErrors.New(status.NoSuchEvent, err)
	} else if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	event.EventHeader = eventHeader

	scheduleRows, err := env.DB().QueryContext(ctx, scheduleQuery, uuid)
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	defer scheduleRows.Close()

	event.EventHeader.Schedule, err = scanObjects(scheduleRows, ScanSchedule)
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	linkRows, err := env.DB().QueryContext(ctx, linksQuery, uuid)
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	defer linkRows.Close()

	event.Links, err = scanObjects(linkRows, ScanLink)
	if err != nil {
		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	return event, nil
}

func scanObjects[T any](rows *sql.Rows, f func(scanner Scanner) (*T, error)) ([]*T, error) {
	objects := make([]*T, 0)
	for rows.Next() {
		object, err := f(rows)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}

	return objects, nil
}
