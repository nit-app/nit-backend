package responses

type ErrorResponse struct {
	BaseResponse[*struct{}]

	Code string `json:"code"`
}
