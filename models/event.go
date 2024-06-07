package models

import "time"

type EventHeader struct {
	UUID               string           `json:"uuid"`
	Title              string           `json:"title" binding:"required"`
	PriceLow           int              `json:"priceLow" binding:"min=0,max=1000000"`
	PriceHigh          int              `json:"priceHigh" binding:"min=0,max=1000000"`
	AgeLimitLow        int              `json:"ageLimitLow" binding:"min=0,max=18"`
	AgeLimitHigh       int              `json:"ageLimitHigh" binding:"min=0,max=99"`
	Location           string           `json:"location" binding:"required"`
	OwnerInfo          string           `json:"ownerInfo" binding:"required"`
	HasCertificate     bool             `json:"hasCertificate"`
	FavCount           int              `json:"favCount"`
	CreatedAt          time.Time        `json:"createdAt"`
	ModifiedAt         time.Time        `json:"modifiedAt"`
	Schedule           []*EventSchedule `json:"schedule"`
	Tags               []string         `json:"tags" binding:"unique,min=0,max=30,dive,min=1,max=512,excludesall= "`
	PlainDescription   string           `json:"plainDescription" binding:"required"`
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
	Title    string    `json:"title" binding:"required,min=2,max=64"`
	URL      string    `json:"url" binding:"required,http_url"`
	AddedAt  time.Time `json:"addedAt"`
}

type EventSchedule struct {
	ScheduleUUID string    `json:"scheduleUUID"`
	BeginsAt     time.Time `json:"beginsAt"`
	EndsAt       time.Time `json:"endsAt"`
	AddedAt      time.Time `json:"addedAt"`
}
