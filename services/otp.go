package services

import (
	"errors"
	"github.com/nit-app/nit-backend/services/otp"
	"github.com/nit-app/nit-backend/services/sms"
	"github.com/nit-app/nit-backend/sessions"
	"github.com/nit-app/nit-backend/validators"
)

type OtpService struct {
	Generator otp.Generator
	Carrier   sms.Carrier
}

func (os *OtpService) Send(session *sessions.Session, phoneNumber string, nextState string) error {
	if session.State != sessions.StateUnauthorized {
		return errBadSignInState
	}

	if err := validators.PhoneNumber(phoneNumber); err != nil {
		return err
	}

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

		return errors.New("otp attempts exceeded")
	}

	if session.OTP.Code != otpCode {
		session.OTP.Attempt++
		return errors.New("bad otp code")
	}

	session.State = nextState

	return nil
}
