package status

import "net/http"

const (
	BadSignInState      = "BAD_SIGN_IN_STATE"
	BadRegisterState    = "BAD_REGISTER_STATE"
	Unauthorized        = "UNAUTHORIZED"
	InvalidDataFormat   = "INVALID_DATA_FORMAT"
	BadRegistrationData = "BAD_REGISTRATION_DATA"
	OtpDeliveryError    = "OTP_DELIVERY_ERROR"
	OtpCheckingError    = "OTP_CHECKING_ERROR"
)

var Codes = map[string]int{Unauthorized: http.StatusUnauthorized}
