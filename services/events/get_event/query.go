package get_event

import (
	"context"
	"github.com/nit-app/nit-backend/env"
	"github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/services/events"
)

func Event(ctx context.Context, uuid *requests.EventIdRequest) (*responses.Event, error) {
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
			e.plainDescription
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

	const descriptionQuery = `
		select
    		e.description
		from 
		    events e
		where 
		    e."uuid" = $1
		    and e.deletedat is null`

	headerRows, err := env.DB().QueryContext(ctx, headerQuery, uuid.UUID)
	if err != nil {
		return nil, errors.New(status.InternalServerError, err)
	}

	defer headerRows.Close()

	eventHeader := &responses.EventHeader{}
	for headerRows.Next() {
		headerObject, err := events.ScanEventHeader(headerRows)
		if err != nil {
			return nil, errors.New(status.InternalServerError, err)
		}
		eventHeader = headerObject
	}

	if eventHeader.UUID == "" {
		return nil, nil
	}

	scheduleRows, err := env.DB().QueryContext(ctx, scheduleQuery, uuid.UUID)
	if err != nil {
		return nil, errors.New(status.InternalServerError, err)
	}

	defer scheduleRows.Close()

	eventSchedules := make([]responses.EventSchedule, 0)
	for scheduleRows.Next() {
		scheduleObject, err := events.ScanSchedule(scheduleRows)
		if err != nil {
			return nil, errors.New(status.InternalServerError, err)
		}
		eventSchedules = append(eventSchedules, scheduleObject)
	}

	linkRows, err := env.DB().QueryContext(ctx, linksQuery, uuid.UUID)
	if err != nil {
		return nil, errors.New(status.InternalServerError, err)
	}

	defer linkRows.Close()

	eventLinks := make([]responses.EventExternalLink, 0)
	for linkRows.Next() {
		linkObject, err := events.ScanLink(linkRows)
		if err != nil {
			return nil, errors.New(status.InternalServerError, err)
		}
		eventLinks = append(eventLinks, linkObject)
	}

	descriptionRows, err := env.DB().QueryContext(ctx, descriptionQuery, uuid.UUID)
	if err != nil {
		return nil, errors.New(status.InternalServerError, err)
	}

	defer descriptionRows.Close()

	var eventDescription string
	for descriptionRows.Next() {
		descriptionObject, err := events.ScanDescription(descriptionRows)
		if err != nil {
			return nil, errors.New(status.InternalServerError, err)
		}
		eventDescription = descriptionObject
	}

	event := &responses.Event{
		EventHeader: *eventHeader,
	}

	event.EventHeader.Schedule = eventSchedules
	event.Links = eventLinks
	event.Description = eventDescription

	return event, nil
}
