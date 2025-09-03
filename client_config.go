package main

import "os"

type ClientConfig struct {
	Broker   string
	Username string
	Password string
	ClientID string
}

func NewClientConfig() ClientConfig {
	broker = os.Getenv("BROKER_URL")
	password = os.Getenv("PASSWORD")
	username = os.Getenv("USERNAME")
	clientID = os.Getenv("CLIENT_ID")

	cConfig := ClientConfig{Broker: broker,
		Password: password,
		Username: username,
		ClientID: clientID}

	return cConfig
}
