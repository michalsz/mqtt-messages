package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/airbrake/gobrake/v5"
	"github.com/michalsz/mqtt_example/clients"
	"github.com/michalsz/mqtt_example/messages"
	"github.com/michalsz/mqtt_example/services"
)

type MessageHandler struct {
	Airbrake *gobrake.Notifier
	Service  services.Sender
}

const (
	msgQueryKey = "msg"
)

func (th MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get(msgQueryKey)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := th.Service.SendMessage(ctx, msg)
	if err != nil {
		th.Airbrake.Notify(err.Error(), nil)
	}
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

type JSONMessageHandler struct {
	Airbrake      *gobrake.Notifier
	Service       services.Sender
	PersistClient clients.PersisterClient
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

		err := th.Service.SendMessage(ctx, dMsg.Value)
		if err != nil {
			th.Airbrake.Notify(err.Error(), nil)
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		addedRecords, err := th.PersistClient.SaveDeviceDatadMsg(dMsg)

		if err != nil {
			th.Airbrake.Notify(err, nil)
			log.Println("error on save records")
		}

		logMsg := fmt.Sprintf("Temp from device: %s Added %d records \n", dMsg.Value, len(addedRecords.Records))

		w.Write([]byte(logMsg))

	} else {
		th.Airbrake.Notify(err, nil)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
