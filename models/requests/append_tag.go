package requests

type AppendTag struct {
	EventUUID string `json:"eventUuid" binding:"required,uuid"`
	Tag       string `json:"tag" binding:"required,min=2,max=15"`
}
