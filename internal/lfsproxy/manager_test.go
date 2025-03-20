package lfsproxy_test

import (
	"context"
	"testing"

	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewProxyManager(t *testing.T) {
	// Test with valid context
	ctx := context.Background()
	pm := lfsproxy.NewProxyManager(ctx)
	assert.NotNil(t, pm)

	// Test with nil context
	pm = lfsproxy.NewProxyManager(nil)
	assert.NotNil(t, pm)

	// Verify empty handlers map is initialized
	assert.NotNil(t, pm)
}

func TestRegisterProxyHandler(t *testing.T) {
	ctx := context.Background()
	pm := lfsproxy.NewProxyManager(ctx)

	// Test registering first handler
	mockHandler := &mockProxyHandler{}
	pm.RegisterProxyHandler("test", mockHandler)

	// Test overwriting existing handler
	mockHandler2 := &mockProxyHandler{}
	pm.RegisterProxyHandler("test", mockHandler2)

	// Test registering nil handler
	pm.RegisterProxyHandler("nil", nil)

	// Test registering handler with empty name
	mockHandler3 := &mockProxyHandler{}
	pm.RegisterProxyHandler("", mockHandler3)

	// Test registering multiple handlers
	pm.RegisterProxyHandler("test1", &mockProxyHandler{})
	pm.RegisterProxyHandler("test2", &mockProxyHandler{})
}

type mockProxyHandler struct{}

func (m *mockProxyHandler) Init(ctx context.Context, r *mux.Router)     {}
func (m *mockProxyHandler) NotifyInitialized()                          {}
