package clients

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const port = 8883
const protocol = "ssl"

func CreateMqttClient(clientConfig ClientConfig) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(ConnectAddress(clientConfig.Broker))
	opts.SetUsername(clientConfig.Username)
	opts.SetPassword(clientConfig.Password)
	opts.SetClientID(clientConfig.ClientID)
	opts.SetKeepAlive(time.Second * 60)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	// if connection failed, exit
	if token.WaitTimeout(3*time.Second) && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return client
}

func ConnectAddress(broker string) string {
	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, broker, port)

	fmt.Println("connect address: ", connectAddress)

	return connectAddress
}
