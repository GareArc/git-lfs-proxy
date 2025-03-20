package lfsproxy

import (
	"context"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy/google"
	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy/lfs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	google_oauth "golang.org/x/oauth2/google"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "",
		Scopes:       []string{"https://www.googleapis.com/auth/drive.appdata"},
		Endpoint:     google_oauth.Endpoint,
	}
)

type GoogleDriveLFSProxy struct {
	ctx        context.Context
	name       string
	lfsManager *lfs.BasicLFSManager
}

func init() {
	RegisterProxyHandler("google", NewGoogleDriveProxy)
}

func NewGoogleDriveProxy() ProxyHandler {
	return &GoogleDriveLFSProxy{
		ctx:  nil,
		name: "google",
	}
}

func (gp *GoogleDriveLFSProxy) Init(ctx context.Context, r *mux.Router) {
	gp.ctx = ctx
	oauthManager := NewOAuthManager(gp.name, oauthConfig, ctx)
	client, err := google.NewGoogleDriveClient(oauthManager.GetTokenSource(), gp.ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Google Drive client")
	}
	lfsHandler := google.NewGoogleDriveLFSHandler(client, gp.ctx, gp.name)
	gp.lfsManager = lfs.NewBasicLFSManager(gp.name, lfsHandler, r, gp.ctx)
}

func (gp *GoogleDriveLFSProxy) NotifyInitialized() {
	log.Info().Str("url", gp.lfsManager.GetLFSFullUrl()).Msg("Google Drive LFS Proxy initialized")
}
