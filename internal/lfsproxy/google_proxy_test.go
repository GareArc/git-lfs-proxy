package lfsproxy_test

import (
	"context"
	"testing"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize required config for tests
	config.Global = &config.Config{
		BaseApiUrl: "http://localhost:8080",
	}
}

func TestNewGoogleDriveProxy(t *testing.T) {
	proxy := lfsproxy.NewGoogleDriveProxy()
	assert.NotNil(t, proxy)
	assert.IsType(t, &lfsproxy.GoogleDriveLFSProxy{}, proxy)
}

func TestNotifyInitialized(t *testing.T) {
	// Skip this test since it requires OAuth setup
	t.Skip("Skipping test that requires OAuth setup")

	proxy := lfsproxy.NewGoogleDriveProxy()
	ctx := context.Background()
	router := mux.NewRouter()

	// Capture log output
	var logOutput string
	log.Logger = log.Logger.Output(testWriter{&logOutput})

	proxy.Init(ctx, router)
	proxy.NotifyInitialized()

	assert.Contains(t, logOutput, "Google Drive LFS Proxy initialized")
}

type testWriter struct {
	output *string
}

func (w testWriter) Write(p []byte) (n int, err error) {
	*w.output += string(p)
	return len(p), nil
}
