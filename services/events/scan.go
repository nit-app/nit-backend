package events

import (
	"github.com/nit-app/nit-backend/models/responses"
	"strings"
)

func ScanEventHeader(row Scanner) (*responses.EventHeader, error) {
	header := &responses.EventHeader{}

	var (
		tags string

		// event header stores only one of the scheduled days, specifically the one that has matched user's query
		matchedDay = responses.EventSchedule{}
	)

	err := row.Scan(&header.UUID, &header.Title, &header.PriceLow, &header.PriceHigh, &header.AgeLimitLow,
		&header.AgeLimitHigh, &header.Location, &header.OwnerInfo, &tags, &header.CreatedAt,
		&header.ModifiedAt, &matchedDay.BeginsAt, &matchedDay.EndsAt, &matchedDay.AddedAt,
		&matchedDay.ScheduleUUID)
	if err != nil {
		return nil, err
	}

	header.Tags = strings.Split(tags, ",")
	header.Schedule = append(header.Schedule, matchedDay)
	return header, nil
}
