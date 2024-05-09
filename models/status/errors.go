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
	TagAlreadySet       = "TAG_ALREADY_SET"
	NoSuchEvent         = "NO_SUCH_EVENT"
	NoSuchUser          = "NO_SUCH_USER"
	Forbidden           = "FORBIDDEN"
)

var Codes = map[string]Code{
	Unauthorized:        {http.StatusUnauthorized, true},
	OtpCheckingError:    {http.StatusUnauthorized, true},
	InternalServerError: {http.StatusInternalServerError, false},
	NoSuchEvent:         {http.StatusNotFound, false},
	NoSuchUser:          {http.StatusNotFound, false},
	Forbidden:           {http.StatusForbidden, false},
	TagAlreadySet:       {http.StatusBadRequest, true},
}

type Code struct {
	HTTP       int
	ExposeText bool
}
