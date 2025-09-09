package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/airbrake/gobrake/v5"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/michalsz/mqtt_example/clients"
	"github.com/michalsz/mqtt_example/handlers"
	"github.com/michalsz/mqtt_example/services"
)

var topic string
var Environment string
var Airbrake *gobrake.Notifier
var AirTblCLient *clients.AirTableCLient

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

	clientConfig := clients.NewClientConfig()
	client := clients.CreateMqttClient(clientConfig)
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
