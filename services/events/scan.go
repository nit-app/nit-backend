package events

import (
	"github.com/nit-app/nit-backend/models/responses"
	"strings"
)

func ScanEventHeader(row Scanner, description *string) (*responses.EventHeader, error) {
	header := &responses.EventHeader{}

	var (
		tags string

		// event header stores only one of the scheduled days, specifically the one that has matched user's query
		matchedDay = &responses.EventSchedule{}
	)

	targets := []any{&header.UUID, &header.Title, &header.PriceLow, &header.PriceHigh, &header.AgeLimitLow,
		&header.AgeLimitHigh, &header.Location, &header.OwnerInfo, &tags, &header.CreatedAt,
		&header.ModifiedAt, &matchedDay.BeginsAt, &matchedDay.EndsAt, &matchedDay.AddedAt,
		&matchedDay.ScheduleUUID, &header.PlainDescription, &header.FavCount}

	if description != nil {
		targets = append(targets, description)
	}

	err := row.Scan(targets...)

	if err != nil {
		return nil, err
	}

	header.Tags = strings.Split(tags, ",")
	header.Schedule = append(header.Schedule, matchedDay)
	return header, nil
}

func ScanSchedule(row Scanner) (*responses.EventSchedule, error) {
	schedule := &responses.EventSchedule{}

	err := row.Scan(&schedule.BeginsAt, &schedule.EndsAt, &schedule.AddedAt, &schedule.ScheduleUUID)

	return schedule, err
}

func ScanLink(row Scanner) (*responses.EventExternalLink, error) {
	link := &responses.EventExternalLink{}

	err := row.Scan(&link.LinkUUID, &link.Title, &link.URL, &link.AddedAt)

	return link, err
}
