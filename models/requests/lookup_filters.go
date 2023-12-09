package requests

import "time"

type EventLookupFilters struct {
	From                 time.Time `json:"from" binding:"required"`
	To                   time.Time `json:"to" binding:"required"`
	Tags                 []string  `json:"tags" binding:"dive,min=1,excludesall= "`
	ExcludeAgeRestricted bool      `json:"excludeAgeRestricted"`
	ExcludePaid          bool      `json:"excludePaid"`
}
