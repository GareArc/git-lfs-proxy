package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create handler and serve request
	healthCheck(rr, req)

	// Check response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check response body
	assert.Equal(t, "OK", rr.Body.String())
}

func TestRegisterHealthCheck(t *testing.T) {
	// Create new router
	router := mux.NewRouter()

	// Register health check endpoint
	registerHealthCheck(router)

	// Create test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Make request to health check endpoint
	resp, err := http.Get(ts.URL + "/health")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
