package messages

import (
	"encoding/json"
	"io"
)

type DeviceMessage struct {
	DeviceId  string `validate:"required"`
	Parameter string `validate:"is_temp"`
	Value     string
	Pressure  float32
	Name      string
}

func DecodeJSON(requestBody io.Reader) (*DeviceMessage, error) {
	dMsg := DeviceMessage{}

	err := json.NewDecoder(requestBody).Decode(&dMsg)
	if err != nil {
		return nil, err
	}

	return &dMsg, nil
}
