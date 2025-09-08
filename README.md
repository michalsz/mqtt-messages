Simple HTTP server with MQTT send message to emqx.com platform.
This is example with one endpoint /send


Run
go run ./...


To send message by http request.
http://localhost:3000/send?msg=new_messaasup


Example json body to endpoint POST /receive

{
    "deviceId": "recWi4wMO4N4fyrbZ",
    "name": "Sensor A2",
    "parameter": "temp",
    "value": "28.3",
    "pressure": 128.3
}