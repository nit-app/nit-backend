package validators

import (
	"github.com/go-playground/validator/v10"
	"strconv"
)

func ruPhoneNumber(fl validator.FieldLevel) bool {
	number := fl.Field().String()
	if len(number) != 11 || number[0] != '7' {
		return false
	}

	_, err := strconv.Atoi(number)
	return err == nil
}
