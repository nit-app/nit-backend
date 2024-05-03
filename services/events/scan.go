package events

import (
	"github.com/nit-app/nit-backend/env"
	"github.com/nit-app/nit-backend/models"
	"strings"
)

func ScanEventHeader(row env.Scanner, description *string) (*models.EventHeader, error) {
	header := &models.EventHeader{}

	var (
		tags string

		// event header stores only one of the scheduled days, specifically the one that has matched user's query
		matchedDay = &models.EventSchedule{}
	)

	targets := []any{&header.UUID, &header.Title, &header.PriceLow, &header.PriceHigh, &header.AgeLimitLow,
		&header.AgeLimitHigh, &header.Location, &header.OwnerInfo, &tags, &header.CreatedAt,
		&header.ModifiedAt, &matchedDay.BeginsAt, &matchedDay.EndsAt, &matchedDay.AddedAt,
		&matchedDay.ScheduleUUID, &header.HasCertificate, &header.PlainDescription, &header.FavCount, &header.IsDraft, &header.IsMachineGenerated}

	if description != nil {
		targets = append(targets, description)
	}

	err := row.Scan(targets...)

	if err != nil {
		return nil, err
	}

	header.Tags = strings.Split(tags, ",")

	header.Schedule = make([]*models.EventSchedule, 0)
	if matchedDay.ScheduleUUID != nil {
		header.Schedule = append(header.Schedule, matchedDay)
	}
	return header, nil
}

func ScanSchedule(row env.Scanner) (*models.EventSchedule, error) {
	schedule := &models.EventSchedule{}

	err := row.Scan(&schedule.BeginsAt, &schedule.EndsAt, &schedule.AddedAt, &schedule.ScheduleUUID)

	return schedule, err
}

func ScanLink(row env.Scanner) (*models.EventExternalLink, error) {
	link := &models.EventExternalLink{}

	err := row.Scan(&link.LinkUUID, &link.Title, &link.URL, &link.AddedAt)

	return link, err
}
