package services

import (
	"database/sql"
	"errors"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/services/user"
	"github.com/nit-app/nit-backend/sessions"
)

type AuthService struct {
	OTP         *OtpService
	UserService *user.Service
}

var errNoUserFoundByNumber = errors.New("no user is created with this phone number")

const maxOtpAttempts = 5

func (as *AuthService) Start(session *sessions.Session, phoneNumber string) error {
	return as.OTP.Send(session, phoneNumber, sessions.StateEnterOtp)
}

func (as *AuthService) CheckOTP(session *sessions.Session, otpCode string) error {
	err := as.OTP.CheckOTP(session, otpCode, sessions.StateEnterOtp)
	if err != nil {
		return err
	}

	subject, err := as.UserService.GetUuidByPhoneNumber(session.OTP.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wrappedErrors.New(status.BadRegistrationData, errNoUserFoundByNumber)
		}

		return err
	}

	sessions.SetAuthorized(session, subject)

	return nil
}
