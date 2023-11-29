package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("ruPhoneNumber", ruPhoneNumber)
		if err != nil {
			panic(err)
		}
	} else {
		panic("cannot register custom validators on engine")
	}
}
