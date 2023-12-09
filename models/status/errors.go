package status

import "net/http"

const (
	BadFormState        = "BAD_FORM_STATE"
	Unauthorized        = "UNAUTHORIZED"
	InvalidDataFormat   = "INVALID_DATA_FORMAT"
	BadRegistrationData = "BAD_REGISTRATION_DATA"
	OtpDeliveryError    = "OTP_DELIVERY_ERROR"
	OtpCheckingError    = "OTP_CHECKING_ERROR"
	InternalServerError = "INTERNAL_SERVER_ERROR"
)

var Codes = map[string]int{Unauthorized: http.StatusUnauthorized, OtpCheckingError: http.StatusUnauthorized,
	InternalServerError: http.StatusInternalServerError}
