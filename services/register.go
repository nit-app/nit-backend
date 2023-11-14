package services

import (
	"github.com/nit-app/nit-backend/models/status"
	"github.com/nit-app/nit-backend/sessions"
)

type RegisterService struct {
	UserService *UserService
	OTP         *OtpService
}

func (as *RegisterService) Start(session *sessions.Session, phoneNumber string) error {
	return as.OTP.Send(session, phoneNumber, sessions.StateRegEnterOtp)
}

func (as *RegisterService) CheckOTP(session *sessions.Session, otpCode string) error {
	return as.OTP.CheckOTP(session, otpCode, sessions.StateRegEnterOtp, sessions.StateRegFinish)
}

func (as *RegisterService) Finish(session *sessions.Session, firstName string, lastName *string) (string, error) {
	newUserUuid, err := as.UserService.RegisterByPhoneNumber(session.OTP.PhoneNumber, firstName, lastName)
	if err != nil {
		return "", err
	}

	sessions.SetAuthorized(session, newUserUuid)

	return newUserUuid, nil
}

func (as *RegisterService) OtpCheckState() string {
	return sessions.StateRegEnterOtp
}

func (as *RegisterService) BadStateStatus() string {
	return status.BadRegisterState
}
