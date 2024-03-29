package services

import (
	stdErrors "errors"
	"github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/services/otp"
	"github.com/nit-app/nit-backend/services/sms"
	"github.com/nit-app/nit-backend/sessions"
)

type OtpService struct {
	Generator otp.Generator
	Carrier   sms.Carrier
}

var (
	errBadOtpState         = errors.New(status.BadFormState, stdErrors.New("bad otp state"))
	errOtpAttemptsExceeded = errors.New(status.OtpCheckingError, stdErrors.New("otp attempts exceeded"))
	errBadOtpCode          = errors.New(status.OtpCheckingError, stdErrors.New("bad otp code"))
)

func (os *OtpService) Send(session *sessions.Session, phoneNumber string, nextState string) error {
	otpCode := os.Generator.Generate()
	session.OTP = &sessions.OtpState{
		Code:        otpCode,
		Attempt:     0,
		PhoneNumber: phoneNumber,
	}
	session.State = nextState
	session.Save()

	return os.Carrier.Send(phoneNumber, otpCode)
}

func (os *OtpService) CheckOTP(session *sessions.Session, otpCode string, expectedState string, nextState string) error {
	if session.State != expectedState || session.OTP == nil {
		return errBadOtpState
	}

	defer session.Save()

	if session.OTP.Attempt >= maxOtpAttempts {
		session.State = sessions.StateUnauthorized
		session.OTP = nil

		return errOtpAttemptsExceeded
	}

	if session.OTP.Code != otpCode {
		session.OTP.Attempt++
		return errBadOtpCode
	}

	session.State = nextState

	return nil
}
