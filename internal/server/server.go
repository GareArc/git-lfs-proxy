package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Router *mux.Router
	Ctx    context.Context
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", config.Global.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	srvErrChan := make(chan error)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			srvErrChan <- err
		}
	}()

	select {
	case err := <-srvErrChan:
		if srvErr := srv.Shutdown(s.Ctx); err != nil {
			log.Error().Err(srvErr).Msg("Failed to shutdown server")
		}
		return err
	case <-s.Ctx.Done():
		log.Info().Msg("Closing server...")
		if err := srv.Shutdown(s.Ctx); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown server")
		}
		return nil
	}
}

func NewServer() (*Server, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	r := mux.NewRouter()
	registerHealthCheck(r)

	proxyManager := lfsproxy.NewProxyManager(ctx)
	proxyManager.InitAllProxies(r)

	//print requests that doesn't match any route
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Warn().Str("path", r.URL.Path).Str("method", r.Method).Msg("Request not found")
		http.Error(w, "Not found", http.StatusNotFound)
	})

	// log registered routes
	if config.Global.LogLevel == "debug" {
		r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			tmpl, err := route.GetPathTemplate()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get path template")
			}
			methods, err := route.GetMethods()
			if err != nil {
				methods = []string{"No method"}
			}
			log.Debug().Str("path", tmpl).Str("method", fmt.Sprintf("%v", methods)).Msg("Registered route")
			return nil
		})
	}

	server := &Server{
		Router: r,
		Ctx:    ctx,
	}
	return server, cancel
}

func registerHealthCheck(r *mux.Router) {
	r.HandleFunc("/health", healthCheck).Methods("GET")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
