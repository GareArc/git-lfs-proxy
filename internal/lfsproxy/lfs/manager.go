package lfs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var (
	LFS_HEADER = "application/vnd.git-lfs+json"
)

/** Basic LFS Manager**/
type BasicLFSHandler interface {
	HandleBatchAPI(w http.ResponseWriter, r *http.Request)
	HandleDownloadAPI(w http.ResponseWriter, r *http.Request)
	HandleUploadAPI(w http.ResponseWriter, r *http.Request)
	HandleVerifyAPI(w http.ResponseWriter, r *http.Request)
	HandleLockAPI(w http.ResponseWriter, r *http.Request)
	HandleUnlockAPI(w http.ResponseWriter, r *http.Request)
	HandleVerifyLockAPI(w http.ResponseWriter, r *http.Request)
}

type BasicLFSManager struct {
	proxyName string
	handler   BasicLFSHandler
	ctx       context.Context
}

func NewBasicLFSManager(proxyName string, handler BasicLFSHandler, r *mux.Router, ctx context.Context) *BasicLFSManager {
	ret := &BasicLFSManager{
		proxyName: proxyName,
		handler:   handler,
		ctx:       ctx,
	}

	ret.registerRouters(r)
	return ret
}

func (m *BasicLFSManager) registerRouters(r *mux.Router) {
	g := r.PathPrefix(m.getApiUrlBase()).Subrouter()
	// g.Use(verifyHeader)
	g.HandleFunc("/objects/batch", m.BatchAPIRouter).Methods("POST")
	g.HandleFunc("/objects/{oid}", m.DownloadAPIRouter).Methods("GET")
	g.HandleFunc("/objects/{oid}", m.UploadAPIRouter).Methods("PUT")
	g.HandleFunc("/objects/{oid}/verify", m.VerifyAPIRouter).Methods("POST")
	g.HandleFunc("/locks/verify", m.VerifyLockAPIRouter).Methods("POST")
	g.HandleFunc("/locks/{oid}", m.LockAPIRouter).Methods("POST")
	g.HandleFunc("/locks/{oid}", m.UnlockAPIRouter).Methods("DELETE")
}

func verifyHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errMsg := ""
		if r.Header.Get("Accept") != LFS_HEADER {
			errMsg = "Invalid Accept header"
		}

		if errMsg != "" {
			log.Error().Msg(errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (m *BasicLFSManager) GetLFSFullUrl() string {
	return fmt.Sprintf("%s%s", config.Global.BaseApiUrl, m.getApiUrlBase())
}

func (m *BasicLFSManager) BatchAPIRouter(w http.ResponseWriter, r *http.Request) {
	m.handler.HandleBatchAPI(w, r)
}

func (m *BasicLFSManager) DownloadAPIRouter(w http.ResponseWriter, r *http.Request) {
	m.handler.HandleDownloadAPI(w, r)
}

func (m *BasicLFSManager) UploadAPIRouter(w http.ResponseWriter, r *http.Request) {
	m.handler.HandleUploadAPI(w, r)
}

func (m *BasicLFSManager) VerifyAPIRouter(w http.ResponseWriter, r *http.Request) {
	m.handler.HandleVerifyAPI(w, r)
}

func (m *BasicLFSManager) LockAPIRouter(w http.ResponseWriter, r *http.Request) {
	m.handler.HandleLockAPI(w, r)
}

func (m *BasicLFSManager) UnlockAPIRouter(w http.ResponseWriter, r *http.Request) {
	m.handler.HandleUnlockAPI(w, r)
}

func (m *BasicLFSManager) VerifyLockAPIRouter(w http.ResponseWriter, r *http.Request) {
	m.handler.HandleVerifyLockAPI(w, r)
}

func (m *BasicLFSManager) getApiUrlBase() string {
	return fmt.Sprintf("/%s/{user}/{repo}/lfs", m.proxyName)
}

/** END Basic LFS Manager**/

type LFSObject struct {
	Oid  string `json:"oid"`
	Size int    `json:"size"`
}

type LFSRef struct {
	Name string `json:"name"`
}

/** BATCH API **/
type LFSBatchRequest struct {
	Operation     string      `json:"operation"`
	Transfers     []string    `json:"transfers,omitempty"`
	Ref           LFSRef      `json:"ref,omitempty"`
	Objects       []LFSObject `json:"objects"`
	HashAlgorithm string      `json:"hash_algo"`
}

type LFSBatchResponse struct {
	Transfer string `json:"transfer"`
	Objects  []any  `json:"objects"`
}

type LFSBatchResponseObject struct {
	Oid           string                            `json:"oid"`
	Size          int                               `json:"size"`
	Authenticated bool                              `json:"authenticated"`
	HashAlgorithm string                            `json:"hash_algo"`
	Actions       map[string]LFSBatchResponseAction `json:"actions,omitempty"`
}

// action types
type LFSBatchResponseAction struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresIn int               `json:"expires_in,omitempty"`
	ExpiresAt string            `json:"expires_at,omitempty"`
}

/** END BATCH API **/
