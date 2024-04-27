package models

import "time"

type EventHeader struct {
	UUID               string           `json:"uuid"`
	Title              string           `json:"title" binding:"required"`
	PriceLow           int              `json:"priceLow" binding:"required"`
	PriceHigh          int              `json:"priceHigh" binding:"required"` // TODO: add validations for remaining fields
	AgeLimitLow        int              `json:"ageLimitLow"`
	AgeLimitHigh       int              `json:"ageLimitHigh"`
	Location           string           `json:"location"`
	OwnerInfo          string           `json:"ownerInfo"`
	HasCertificate     bool             `json:"hasCertificate"`
	FavCount           int              `json:"favCount"`
	CreatedAt          time.Time        `json:"createdAt"`
	ModifiedAt         time.Time        `json:"modifiedAt"`
	Schedule           []*EventSchedule `json:"schedule"`
	Tags               []string         `json:"tags"`
	PlainDescription   string           `json:"plainDescription"`
	IsDraft            bool             `json:"isDraft"`
	IsMachineGenerated bool             `json:"isMachineGenerated"`
}

type Event struct {
	*EventHeader
	Links       []*EventExternalLink `json:"links"`
	Description string               `json:"description"`
	// photos, markdown, reviews, etc..
}

type EventExternalLink struct {
	LinkUUID string    `json:"uuid"`
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	AddedAt  time.Time `json:"addedAt"`
}

type EventSchedule struct {
	ScheduleUUID string    `json:"scheduleUUID"`
	BeginsAt     time.Time `json:"beginsAt"`
	EndsAt       time.Time `json:"endsAt"`
	AddedAt      time.Time `json:"addedAt"`
}
