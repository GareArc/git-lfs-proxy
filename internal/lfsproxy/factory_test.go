package lfsproxy

import (
	"context"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Mock ProxyHandler for testing
type MockProxyHandler struct{}

func (m *MockProxyHandler) Init(ctx context.Context, r *mux.Router)    {}
func (m *MockProxyHandler) NotifyInitialized()                         {}

func TestRegisterProxyHandler(t *testing.T) {
	// Reset map for testing
	proxyFactoryMap = make(map[string]func() ProxyHandler)

	mockHandler := &MockProxyHandler{}
	mockBuilder := func() ProxyHandler {
		return mockHandler
	}

	// Test registering new handler
	RegisterProxyHandler("test", mockBuilder)
	assert.Equal(t, 1, len(proxyFactoryMap), "Expected map length 1")

	// Test overwriting existing handler
	newMockHandler := &MockProxyHandler{}
	newMockBuilder := func() ProxyHandler {
		return newMockHandler
	}
	RegisterProxyHandler("test", newMockBuilder)
	assert.Equal(t, 1, len(proxyFactoryMap), "Expected map length 1")

	// Test registering empty name
	RegisterProxyHandler("", mockBuilder)
	assert.Equal(t, 2, len(proxyFactoryMap), "Expected map length 2")

	// Test registering nil builder
	RegisterProxyHandler("nil", nil)
	assert.Equal(t, 3, len(proxyFactoryMap), "Expected map length 3")
}

func TestGetProxyHandler(t *testing.T) {
	// Reset map for testing
	proxyFactoryMap = make(map[string]func() ProxyHandler)

	mockHandler := &MockProxyHandler{}
	mockBuilder := func() ProxyHandler {
		return mockHandler
	}

	// Test getting non-existent handler
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-existent handler")
		}
	}()

	// Register and get handler
	RegisterProxyHandler("test", mockBuilder)
	handler := GetProxyHandler("test")

	// Type assertion to compare concrete types
	mockHandlerResult, ok := handler.(*MockProxyHandler)
	assert.True(t, ok, "Expected handler to be *MockProxyHandler")
	assert.Equal(t, mockHandler, mockHandlerResult)

	// This should panic
	GetProxyHandler("nonexistent")
}
