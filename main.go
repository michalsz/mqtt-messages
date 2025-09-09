package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/airbrake/gobrake/v5"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/davecgh/go-spew/spew"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/michalsz/mqtt_example/clients"
	"github.com/michalsz/mqtt_example/handlers"
	"github.com/michalsz/mqtt_example/services"
)

var broker string
var password string
var username string
var topic string
var clientID string
var Environment string
var Airbrake *gobrake.Notifier
var AirTblCLient *clients.AirTableCLient

const port = 8883
const protocol = "ssl"

func createMqttClient(clientConfig ClientConfig) mqtt.Client {
	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, clientConfig.Broker, port)

	fmt.Println("connect address: ", connectAddress)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectAddress)
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	projectID, _ := strconv.ParseInt(os.Getenv("AIRBRAKE_PROJECT_ID"), 10, 64)
	projectKey := os.Getenv("AIRBRAKE_PROJECT_KEY")
	Environment := os.Getenv("ENVIRONMENT")

	spew.Dump(Environment)

	Airbrake = gobrake.NewNotifierWithOptions(&gobrake.NotifierOptions{
		ProjectId:   projectID,
		ProjectKey:  projectKey,
		Environment: Environment,
	})
}

func estamblishTopic() {
	flag.StringVar(&topic, "topic", "t/1", "a string")
	flag.Parse()
}

func main() {
	defer Airbrake.Close()
	defer Airbrake.NotifyOnPanic()

	clientConfig := NewClientConfig()
	client := createMqttClient(clientConfig)
	estamblishTopic()

	mux := http.NewServeMux()

	service := services.QtMessageSender{Client: client}

	msgHandler := handlers.MessageHandler{
		Airbrake: Airbrake,
		Service:  service,
	}
	mux.Handle("/send", msgHandler)

	c := clients.NewAirTableClient()

	jsonHandler := handlers.JSONMessageHandler{
		Airbrake:      Airbrake,
		Service:       service,
		PersistClient: c,
	}
	mux.Handle("/receive", jsonHandler)
	mux.HandleFunc("POST /receive/", jsonHandler.ServeHTTP)

	hHandler := handlers.HealthCheckHandler{}
	mux.Handle("/health", hHandler)

	Environment := os.Getenv("ENVIRONMENT")
	spew.Dump(Environment)
	if Environment == "development" {
		http.ListenAndServe(":3000", mux)
	} else {
		lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
	}
}
