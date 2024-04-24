package auth

import (
	"github.com/nit-app/nit-backend/services/user"
	"github.com/nit-app/nit-backend/sessions"
)

type RegisterBundle struct{}

var Register = &RegisterBundle{}

func (r *RegisterBundle) Start(session *sessions.Session, phoneNumber string) error {
	return otpService.Send(session, phoneNumber, sessions.StateRegEnterOtp)
}

func (r *RegisterBundle) CheckOTP(session *sessions.Session, otpCode string) error {
	err := otpService.CheckOTP(session, otpCode, sessions.StateRegEnterOtp)
	if err != nil {
		return err
	}

	session.State = sessions.StateRegFinish
	session.Save()

	return nil
}

func (r *RegisterBundle) Finish(session *sessions.Session, firstName string, lastName *string) (string, error) {
	newUserUuid, err := user.RegisterByPhoneNumber(session.OTP.PhoneNumber, firstName, lastName)
	if err != nil {
		return "", err
	}

	sessions.SetAuthorized(session, newUserUuid)

	return newUserUuid, nil
}
