package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/airbrake/gobrake/v5"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-playground/validator/v10"
)

type MessageHandler struct {
	Client   mqtt.Client
	Airbrake *gobrake.Notifier
}

const (
	mqttTopic      = "new_topic"
	mqttQoS        = 0 // Quality of Service: 0 means "at most once".
	testMessageErr = "error"
	msgQueryKey    = "msg"
)

func sendMessage(message string, client mqtt.Client) error {
	payload := fmt.Sprintf("message: %s!", message)
	if token := client.Publish(mqttTopic, byte(mqttQoS), false, payload); token.Wait() && token.Error() != nil {
		log.Printf("publish failed, topic: %s, payload: %s\n", mqttTopic, payload)
	} else {
		log.Printf("publish success, topic: %s, payload: %s\n", mqttTopic, payload)
	}

	if message == testMessageErr {
		return errors.New("error Test from Airbrake")
	}

	return nil
}

func (th MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get(msgQueryKey)
	err := sendMessage(msg, th.Client)
	if err != nil {
		th.Airbrake.Notify(err.Error(), nil)
	}
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

type JSONMessageHandler struct {
	Client   mqtt.Client
	Airbrake *gobrake.Notifier
}

type DeviceMessage struct {
	DeviceId  string `validate:"required"`
	Parameter string `validate:"is_temp"`
	Value     string
}

func (th JSONMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dMsg, err := decodeJSON(r.Body)

	if err != nil {
		th.Airbrake.Notify(err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	isValid, err := validateMsg(dMsg)
	if isValid {
		err = sendMessage(dMsg.Value, th.Client)
		if err != nil {
			th.Airbrake.Notify(err.Error(), nil)
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	} else {
		th.Airbrake.Notify(err, nil)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	w.Write([]byte("Temp from device: " + dMsg.Value))
}

func decodeJSON(requestBody io.Reader) (*DeviceMessage, error) {
	dMsg := DeviceMessage{}

	err := json.NewDecoder(requestBody).Decode(&dMsg)
	if err != nil {
		return nil, err
	}

	return &dMsg, nil
}

func validateMsg(dMsg *DeviceMessage) (bool, error) {
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
