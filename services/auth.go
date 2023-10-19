package services

import (
	"database/sql"
	"errors"
	"github.com/nit-app/nit-backend/sessions"
)

type AuthService struct {
	OTP         *OtpService
	UserService *UserService
}

const maxOtpAttempts = 5

var (
	errBadSignInState = errors.New("bad sign in state")
	errBadOtpState    = errors.New("bad otp state")
)

func (as *AuthService) Start(session *sessions.Session, phoneNumber string) error {
	return as.OTP.Send(session, phoneNumber, sessions.StateEnterOtp)
}

func (as *AuthService) CheckOTP(session *sessions.Session, otpCode string) error {
	err := as.OTP.CheckOTP(session, otpCode, sessions.StateEnterOtp, sessions.StateAuthorized)
	if err != nil {
		return err
	}

	subject, err := as.UserService.GetUuidByPhoneNumber(session.OTP.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("no user is created with this phone number")
		}

		return err
	}

	sessions.SetAuthorized(session, subject)

	return nil
}

func (as *AuthService) OtpCheckState() string {
	return sessions.StateEnterOtp
}
