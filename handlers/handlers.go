package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/airbrake/gobrake/v5"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageHandler struct {
	Client   mqtt.Client
	Airbrake *gobrake.Notifier
}

const (
	mqttTopic      = "new_topic"
	mqttQoS        = 0 // Quality of Service: 0 means "at most once".
	testMessageErr = "error"
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
	msg := r.URL.Query().Get("msg")
	err := sendMessage(msg, th.Client)
	if err != nil {
		th.Airbrake.Notify(err.Error(), nil)
	}
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}
