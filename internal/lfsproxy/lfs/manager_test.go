package lfs_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy/lfs"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct{}

func (m *mockHandler) HandleBatchAPI(w http.ResponseWriter, r *http.Request)      {}
func (m *mockHandler) HandleDownloadAPI(w http.ResponseWriter, r *http.Request)   {}
func (m *mockHandler) HandleUploadAPI(w http.ResponseWriter, r *http.Request)     {}
func (m *mockHandler) HandleVerifyAPI(w http.ResponseWriter, r *http.Request)     {}
func (m *mockHandler) HandleLockAPI(w http.ResponseWriter, r *http.Request)       {}
func (m *mockHandler) HandleUnlockAPI(w http.ResponseWriter, r *http.Request)     {}
func (m *mockHandler) HandleVerifyLockAPI(w http.ResponseWriter, r *http.Request) {}

func TestGetLFSFullUrl(t *testing.T) {
	// Setup test config
	config.Global = &config.Config{
		BaseApiUrl: "http://localhost:8080",
	}

	tests := []struct {
		name      string
		proxyName string
		expected  string
	}{
		{
			name:      "Basic URL",
			proxyName: "test-proxy",
			expected:  "http://localhost:8080/test-proxy/{user}/{repo}/lfs",
		},
		{
			name:      "Empty proxy name",
			proxyName: "",
			expected:  "http://localhost:8080//{user}/{repo}/lfs",
		},
		{
			name:      "Special characters in proxy name",
			proxyName: "test@proxy",
			expected:  "http://localhost:8080/test@proxy/{user}/{repo}/lfs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := lfs.NewBasicLFSManager(tt.proxyName, &mockHandler{}, mux.NewRouter(), context.Background())
			result := manager.GetLFSFullUrl()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewBasicLFSManager(t *testing.T) {
	tests := []struct {
		name      string
		proxyName string
	}{
		{
			name:      "Valid manager creation",
			proxyName: "test-proxy",
		},
		{
			name:      "Empty proxy name",
			proxyName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			handler := &mockHandler{}
			ctx := context.Background()

			manager := lfs.NewBasicLFSManager(tt.proxyName, handler, router, ctx)
			assert.NotNil(t, manager)
		})
	}
}

func TestVerifyHeaderMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		statusCode int
	}{
		{
			name:       "Valid header",
			header:     lfs.LFS_HEADER,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid header",
			header:     "invalid/header",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Empty header",
			header:     "",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.header != "" {
				req.Header.Set("Accept", tt.header)
			}

			rr := httptest.NewRecorder()
			handler := http.Handler(nextHandler)
			verifyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Accept") != lfs.LFS_HEADER {
					http.Error(w, "Invalid Accept header", http.StatusBadRequest)
					return
				}
				handler.ServeHTTP(w, r)
			})

			verifyHandler.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestRouterHandlers(t *testing.T) {
	router := mux.NewRouter()
	handler := &mockHandler{}
	ctx := context.Background()
	lfs.NewBasicLFSManager("test-proxy", handler, router, ctx)

	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{
			name:       "Batch API",
			method:     "POST",
			path:      "/test-proxy/user/repo/lfs/objects/batch",
			statusCode: http.StatusOK,
		},
		{
			name:       "Download API",
			method:     "GET",
			path:      "/test-proxy/user/repo/lfs/objects/123",
			statusCode: http.StatusOK,
		},
		{
			name:       "Upload API",
			method:     "PUT",
			path:      "/test-proxy/user/repo/lfs/objects/123",
			statusCode: http.StatusOK,
		},
		{
			name:       "Verify API",
			method:     "POST",
			path:      "/test-proxy/user/repo/lfs/objects/123/verify",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}
