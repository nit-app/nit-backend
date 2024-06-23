package status

import "net/http"

const (
	BadFormState             = "BAD_FORM_STATE"
	Unauthorized             = "UNAUTHORIZED"
	InvalidDataFormat        = "INVALID_DATA_FORMAT"
	InvalidDataFormatUnknown = "INVALID_DATA_FORMAT_UNK" // for db-side checks; does not expose actual error via api
	BadRegistrationData      = "BAD_REGISTRATION_DATA"
	OtpDeliveryError         = "OTP_DELIVERY_ERROR"
	OtpCheckingError         = "OTP_CHECKING_ERROR"
	InternalServerError      = "INTERNAL_SERVER_ERROR"
	DuplicateValueEntry      = "DUPLICATE_VALUE_ENTRY"
	MalformedEvent           = "MALFORMED_EVENT"
	AlreadyPublished         = "ALREADY_PUBLISHED"
	NoSuchEvent              = "NO_SUCH_EVENT"
	NoSuchUser               = "NO_SUCH_USER"
	Forbidden                = "FORBIDDEN"
)

type Code struct {
	HTTP       int
	ExposeText bool
}

var Codes = map[string]Code{
	Unauthorized:             {http.StatusUnauthorized, true},
	OtpCheckingError:         {http.StatusUnauthorized, true},
	InternalServerError:      {http.StatusInternalServerError, false},
	NoSuchEvent:              {http.StatusNotFound, false},
	NoSuchUser:               {http.StatusNotFound, false},
	Forbidden:                {http.StatusForbidden, false},
	DuplicateValueEntry:      {http.StatusBadRequest, true},
	InvalidDataFormat:        {http.StatusBadRequest, true},
	InvalidDataFormatUnknown: {http.StatusBadRequest, false},
	MalformedEvent:           {http.StatusBadRequest, true},
}
