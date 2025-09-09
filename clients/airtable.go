package clients

import (
	"os"
	"strconv"
	"time"

	"github.com/mehanizm/airtable"
	"github.com/michalsz/mqtt_example/messages"
)

type PersisterClient interface {
	SaveDeviceDatadMsg(dMsg *messages.DeviceMessage) (*airtable.Records, error)
}

type AirTableCLient struct {
	Client   *airtable.Client
	dbName   string
	viewName string
}

func (atClient AirTableCLient) GetTable() *airtable.Table {
	return atClient.Client.GetTable(atClient.dbName, atClient.viewName)
}

func (atClient AirTableCLient) SaveDeviceDatadMsg(dMsg *messages.DeviceMessage) (*airtable.Records, error) {
	recordsToSend := &airtable.Records{
		Records: []*airtable.Record{
			{
				Fields: prepareDeviceDataRecord(*dMsg),
			},
		},
	}
	records, err := atClient.GetTable().AddRecords(recordsToSend)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func NewAirTableClient() *AirTableCLient {
	airtableToken := os.Getenv("AIRTABLE_TOKEN")
	airtableDBName := os.Getenv("AIRTABLE_DB_NAME")
	airtableViewName := os.Getenv("AIRTABLE_VIEW_NAME")

	airTableClient := airtable.NewClient(airtableToken)

	return &AirTableCLient{Client: airTableClient,
		dbName:   airtableDBName,
		viewName: airtableViewName,
	}
}

func prepareDeviceDataRecord(dMsg messages.DeviceMessage) map[string]any {
	tempVal, _ := strconv.ParseFloat(dMsg.Value, 64)

	deviceData := make(map[string]any)
	deviceData["Name"] = dMsg.Name
	deviceData["Device"] = []string{dMsg.DeviceId}
	deviceData["Temperature"] = tempVal
	deviceData["Timestamp"] = time.Now().Format("01.02.2006")
	deviceData["Pressure"] = dMsg.Pressure

	return deviceData
}

// records, err := table.GetRecords().
// 	FromView("all_data").
// 	// ReturnFields("Name").
// 	Do()

// if err != nil {
// 	log.Println("error on get records")
// 	log.Println(err)
// }
// for _, r := range records.Records {
// 	log.Println(r.ID)
// 	log.Println(r.Fields["Name"])
// }

// tempVal, _ := strconv.Atoi(dMsg.Value)
// timestamp := time.Now().Format("01.02.2006")
// dId := []string{dMsg.DeviceId}
