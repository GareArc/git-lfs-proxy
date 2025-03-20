package lfsproxy

import (
	"context"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type ProxyManager struct {
	proxyHandlers map[string]ProxyHandler
	ctx           context.Context
}

func NewProxyManager(ctx context.Context) *ProxyManager {
	ret := &ProxyManager{
		proxyHandlers: make(map[string]ProxyHandler),
		ctx:           ctx,
	}

	return ret
}

func (pm *ProxyManager) RegisterProxyHandler(name string, handler ProxyHandler) {
	pm.proxyHandlers[name] = handler
}

func (pm *ProxyManager) InitAllProxies(r *mux.Router) {

	if config.Global.GoogleDriveProxyConfig.Enabled {
		pm.RegisterProxyHandler("google", GetProxyHandler("google"))
		pm.proxyHandlers["google"].Init(pm.ctx, r)
	}

	// notify all proxies started
	for _, handler := range pm.proxyHandlers {
		handler.NotifyInitialized()
	}

	// handle when no proxy is enabled
	if len(pm.proxyHandlers) == 0 {
		log.Warn().Msg("No proxy is enabled")
	}

}
