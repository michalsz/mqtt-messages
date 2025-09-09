package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mehanizm/airtable"
	"github.com/michalsz/mqtt_example/messages"
)

type mockService struct {
}

func (s mockService) SendMessage(ctx context.Context, message string) error {
	return nil
}

func TestSendHandler(t *testing.T) {
	validData := []byte(``)
	req := httptest.NewRequest(http.MethodGet, "/send?msg=test", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	mService := mockService{}
	h := MessageHandler{Service: mService}
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "The time is") {
		t.Errorf("Expected body %s, got %s", "The time is", w.Body)
	}
}

type mockPersisterClient struct {
}

func (c mockPersisterClient) SaveDeviceDatadMsg(dMsg *messages.DeviceMessage) (*airtable.Records, error) {
	records := airtable.Records{}
	return &records, nil
}

func TestSendHandlerReceive(t *testing.T) {
	validData := []byte(`{"deviceId": "recWi4wMO4N4fyrbZ",
    "name": "Sensor A2",
    "parameter": "temp",
    "value": "23.3",
    "pressure": 118.3}`)
	req := httptest.NewRequest(http.MethodPost, "/receive", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	mService := mockService{}
	mCLient := mockPersisterClient{}
	h := JSONMessageHandler{Service: mService, PersistClient: mCLient}
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Temp from device: 23.3 Added") {
		t.Errorf("Expected body %s, got %s", "Temp from device: 23.3 Added", w.Body)
	}
}
