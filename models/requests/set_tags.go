package requests

type SetTags struct {
	EventUUID string   `json:"eventUuid" binding:"required,uuid"`
	Tags      []string `json:"tags" binding:"required,unique,min=1,max=30,dive,min=3,max=512,excludesall= "`
}
