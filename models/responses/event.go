package responses

import "time"

type EventHeader struct {
	UUID             string           `json:"uuid"`
	Title            string           `json:"title"`
	PriceLow         int              `json:"priceLow"`
	PriceHigh        int              `json:"priceHigh"`
	AgeLimitLow      int              `json:"ageLimitLow"`
	AgeLimitHigh     int              `json:"ageLimitHigh"`
	Location         string           `json:"location"`
	OwnerInfo        string           `json:"ownerInfo"`
	HasCertificate   bool             `json:"hasCertificate"`
	FavCount         int              `json:"favCount"`
	CreatedAt        time.Time        `json:"createdAt"`
	ModifiedAt       time.Time        `json:"modifiedAt"`
	Schedule         []*EventSchedule `json:"schedule"`
	Tags             []string         `json:"tags"`
	PlainDescription string           `json:"plainDescription"`
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
