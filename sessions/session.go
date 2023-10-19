package sessions

import "errors"

const (
	StateUnauthorized = "unauthorized"
	StateEnterOtp     = "enter_otp"
	StateRegEnterOtp  = "reg_enter_otp"
	StateRegFinish    = "reg_finish"
	StateAuthorized   = "authorized"
)

const SessionKey = "app_session_current_user"

type Session struct {
	State      string
	Authorized int64

	OTP     *OtpState
	Subject *string

	tokHash string
}

var ErrNoSuchSession = errors.New("no such session")
