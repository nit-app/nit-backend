package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/env"
	"time"
)

type Service struct {
}

var (
	errPhoneNumberOccupied = errors.New("this phone number already has an account associated with it")
)

func (us *Service) GetUuidByPhoneNumber(phoneNumber string) (string, error) {
	row := env.DB().QueryRow("select uuid from users where phoneNumber = $1", phoneNumber)

	var userUuid string
	err := row.Scan(&userUuid)

	return userUuid, err
}

func (us *Service) RegisterByPhoneNumber(phoneNumber string, firstName string, lastName *string) (string, error) {
	existing, _ := us.GetUuidByPhoneNumber(phoneNumber)
	if len(existing) != 0 {
		return "", errPhoneNumberOccupied
	}

	userUuid := uuid.New()
	_, err := env.DB().Exec("insert into users (uuid, phoneNumber, firstName, lastName, registeredAt) values ($1, $2, $3, $4, $5)", userUuid.String(), phoneNumber, firstName, lastName, time.Now())
	if err != nil {
		return "", err
	}

	return userUuid.String(), err
}

func New() *Service {
	return &Service{}
}
