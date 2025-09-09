package clients

import (
	"testing"

	"github.com/michalsz/mqtt_example/clients"
)


func TestConnectAddress(t *testing.T) {
	tests := []struct {
		broker   string
		expected string
	}{
		{"test.mqtt-broker.com", "ssl://test.mqtt-broker.com:8883"},
		{"localhost", "ssl://localhost:8883"},
	}

	// protocol := "ssl"
	// port := 8883
	for _, tt := range tests {
		got := clients.ConnectAddress(tt.broker)
		// got := fmt.Sprintf("%s://%s:%d", protocol, tt.broker, port)
		if got != tt.expected {
			t.Errorf("connectAddress(%q) = %q; want %q", tt.broker, got, tt.expected)
		}
	}
}
