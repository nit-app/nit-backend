package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/env"
	"github.com/nit-app/nit-backend/validators"
	"time"
)

type UserService struct {
}

func (us *UserService) GetUuidByPhoneNumber(phoneNumber string) (string, error) {
	if err := validators.PhoneNumber(phoneNumber); err != nil {
		return "", err
	}

	row := env.DB().QueryRow("select uuid from users where phoneNumber = $1", phoneNumber)

	var userUuid string
	err := row.Scan(&userUuid)

	return userUuid, err
}

func (us *UserService) RegisterByPhoneNumber(phoneNumber string, firstName string, lastName *string) (string, error) {
	existing, _ := us.GetUuidByPhoneNumber(phoneNumber)
	if len(existing) != 0 {
		return "", errors.New("this phone number already has an account associated with it")
	}

	userUuid := uuid.New()
	_, err := env.DB().Exec("insert into users (uuid, phoneNumber, firstName, lastName, registeredAt) values ($1, $2, $3, $4, $5)", userUuid.String(), phoneNumber, firstName, lastName, time.Now())
	if err != nil {
		return "", err
	}

	return userUuid.String(), err
}