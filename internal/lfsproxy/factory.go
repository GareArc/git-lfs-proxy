package lfsproxy

import (
	"context"

	"github.com/gorilla/mux"
)

type ProxyHandler interface {
	Init(ctx context.Context, r *mux.Router)
	NotifyInitialized()
}

var (
	proxyFactoryMap map[string]func() ProxyHandler = make(map[string]func() ProxyHandler)
)

func RegisterProxyHandler(name string, builder func() ProxyHandler) {
	proxyFactoryMap[name] = builder
}

func GetProxyHandler(name string) ProxyHandler {
	return proxyFactoryMap[name]()
}
