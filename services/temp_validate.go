package services

import (
	"github.com/go-playground/validator/v10"
	"github.com/michalsz/mqtt_example/messages"
)

func ValidateMsg(dMsg *messages.DeviceMessage) (bool, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("is_temp", tempValidator)

	err := validate.Struct(dMsg)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func tempValidator(fl validator.FieldLevel) bool {
	tmp := fl.Field().String()
	return tmp == "temp"
}
