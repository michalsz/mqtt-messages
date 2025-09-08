package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/airbrake/gobrake/v5"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/michalsz/mqtt_example/messages"
	"github.com/michalsz/mqtt_example/services"
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

func sendMessage(ctx context.Context, message string, client mqtt.Client) error {
	payload := fmt.Sprintf("message: %s!", message)
	results := make(chan mqtt.Token, 1)

	go func() {
		token := client.Publish(mqttTopic, byte(mqttQoS), false, payload)

		if message == "100" {
			time.Sleep(3 * time.Second)
		}
		if token.Wait() && token.Error() != nil {
			log.Printf("publish failed, topic: %s, payload: %s\n", mqttTopic, payload)
		} else {
			log.Printf("publish success, topic: %s, payload: %s\n", mqttTopic, payload)
		}
		results <- token
	}()

	if message == testMessageErr {
		return errors.New("error Test from Airbrake")
	}

	select {
	case <-results:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (th MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get(msgQueryKey)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := sendMessage(ctx, msg, th.Client)
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

func (th JSONMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dMsg, err := messages.DecodeJSON(r.Body)

	if err != nil {
		th.Airbrake.Notify(err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	isValid, err := services.ValidateMsg(dMsg)
	if isValid {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		err = sendMessage(ctx, dMsg.Value, th.Client)
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
