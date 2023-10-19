package validators

import (
	"errors"
	"strconv"
)

var (
	ErrBadPhoneNumber = errors.New("validation error: bad phone number, must be 7XXXXXXXXXX")
)

func PhoneNumber(number string) error {
	if len(number) != 11 || number[0] != '7' {
		return ErrBadPhoneNumber
	}

	_, err := strconv.Atoi(number)
	return err
}
