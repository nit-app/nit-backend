package auth

import (
	"database/sql"
	"errors"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/services/user"
	"github.com/nit-app/nit-backend/sessions"
)

var errNoUserFoundByNumber = errors.New("no user is created with this phone number")

var otpService = &OtpService{}

const maxOtpAttempts = 5

type SignInBundle struct{}

var SignIn = &SignInBundle{}

func (s *SignInBundle) Start(session *sessions.Session, phoneNumber string) error {
	return otpService.Send(session, phoneNumber, sessions.StateEnterOtp)
}

func (s *SignInBundle) CheckOTP(session *sessions.Session, otpCode string) error {
	err := otpService.CheckOTP(session, otpCode, sessions.StateEnterOtp)
	if err != nil {
		return err
	}

	subject, err := user.GetUuidByPhoneNumber(session.OTP.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wrappedErrors.New(status.BadRegistrationData, errNoUserFoundByNumber)
		}

		return err
	}

	sessions.SetAuthorized(session, subject)

	return nil
}
