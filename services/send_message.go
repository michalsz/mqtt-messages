package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttTopic      = "new_topic"
	mqttQoS        = 0 // Quality of Service: 0 means "at most once".
	testMessageErr = "error"
)

type Sender interface {
	SendMessage(ctx context.Context, message string) error
}

type QtMessageSender struct {
	Client mqtt.Client
}

func (s QtMessageSender) SendMessage(ctx context.Context, message string) error {
	payload := fmt.Sprintf("message: %s!", message)
	results := make(chan mqtt.Token, 1)

	go func() {
		token := s.Client.Publish(mqttTopic, byte(mqttQoS), false, payload)

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
