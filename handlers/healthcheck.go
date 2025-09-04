package handlers

import (
	"net/http"
)

type HealthCheckHandler struct {
}

func (hHandler HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
