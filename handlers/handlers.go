package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/airbrake/gobrake/v5"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/michalsz/mqtt_example/clients"
	"github.com/michalsz/mqtt_example/messages"
	"github.com/michalsz/mqtt_example/services"
)

type MessageHandler struct {
	Client   mqtt.Client
	Airbrake *gobrake.Notifier
}

const (
	msgQueryKey = "msg"
)

func (th MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get(msgQueryKey)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := services.SendMessage(ctx, msg, th.Client)
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
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err = services.SendMessage(ctx, dMsg.Value, th.Client)
		if err != nil {
			th.Airbrake.Notify(err.Error(), nil)
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		c := clients.NewAirTableClient()
		addedRecords, err := c.SaveDeviceDatadMsg(dMsg)

		if err != nil {
			th.Airbrake.Notify(err, nil)
			log.Println("error on save records")
		}

		logMsg := fmt.Sprintf("Temp from device: %s Added %d records \n", dMsg.Value, len(addedRecords.Records))

		log.Println(logMsg)
		w.Write([]byte(logMsg))

	} else {
		th.Airbrake.Notify(err, nil)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
