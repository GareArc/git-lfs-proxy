package lfsproxy

import (
	"context"
	"fmt"
	"net/http"
	url_tool "net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/db"
	"github.com/GareArc/git-lfs-proxy/internal/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type OAuthManager struct {
	proxyOwner  string
	oauthConfig *oauth2.Config
	credential  *models.OAuthCredential
	loginChan   chan any
	ctx         context.Context
}

func NewOAuthManager(proxyOwner string, oauthConfig *oauth2.Config, ctx context.Context) *OAuthManager {
	ret := &OAuthManager{
		oauthConfig: oauthConfig,
		proxyOwner:  proxyOwner,
		credential:  nil,
		loginChan:   make(chan any),
		ctx:         ctx,
	}
	ret.oauthConfig.RedirectURL = ret.getFullOAuthCallbackUrl()

	ret.ensureCredential()
	return ret
}

func (o *OAuthManager) GetCredential() *models.OAuthCredential {
	return o.credential
}

func (o *OAuthManager) GetTokenSource() oauth2.TokenSource {
	return o.oauthConfig.TokenSource(o.ctx, &o.credential.Token)
}

func (o *OAuthManager) ensureCredential() error {
	// try to find existing credential in database
	cred := &models.OAuthCredential{}
	db.DB.First(&cred, "proxy_owner = ?", o.proxyOwner)
	if cred.ID != 0 {
		o.credential = cred
		err := o.refreshOauthCredential()
		if err == nil { // all good, proceed
			return nil
		}
	}

	// credential not found or failed to refresh
	// create temporary http server to handle oauth callback
	srv := o.createTempHTTPServer()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Temporary server failed")
		}
	}()

	o.oauthLoginBrowser()
	log.Info().Msg("Please login to Google Drive first...")
	// wait for login close
	<-o.loginChan
	log.Info().Msg("Login successful")

	// close the temporary http server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		cancel()
	}

	return nil
}

func (o *OAuthManager) oauthLoginBrowser() {
	url := o.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	err := openURL(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open browser")
	}
}

func (o *OAuthManager) createTempHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/%s/oauth/callback", o.proxyOwner), o.handleGoogleOAuthCallback)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Global.Port),
		Handler: mux,
	}

}

func (o *OAuthManager) getFullOAuthCallbackUrl() string {
	return fmt.Sprintf("%s/%s/oauth/callback", config.Global.BaseApiUrl, o.proxyOwner)
}

func (o *OAuthManager) refreshOauthCredential() error {
	if o.credential == nil {
		return fmt.Errorf("credential is nil")
	}
	tokenSource := o.oauthConfig.TokenSource(o.ctx, &o.credential.Token)
	token, err := tokenSource.Token()
	if err != nil {
		return err
	}

	o.credential.Token = *token
	db.DB.Save(o.credential)
	return nil
}

func (o *OAuthManager) handleGoogleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Error().Msg("Code not found")
		w.Write([]byte("Code not found"))
		return
	}

	token, err := o.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Error().Err(err).Msg("Failed to exchange token")
		w.Write([]byte("Failed to exchange token"))
		return
	}

	log.Debug().Msgf("Token: %v", token)

	// find existing credential in database
	cred := &models.OAuthCredential{}
	db.DB.First(&cred, "proxy_owner = ?", "google")
	if cred.ID != 0 {
		cred.Token = *token
		db.DB.Save(&cred)
		o.credential = cred
	} else {
		oauthCred := &models.OAuthCredential{
			ProxyOwner: o.proxyOwner,
			Token:      *token,
		}
		db.DB.Create(oauthCred)
		o.credential = oauthCred
	}

	w.Write([]byte("Succeed. You can close this tab now."))

	if o.loginChan != nil {
		// send signal to loginChan by closing it
		close(o.loginChan)
	}
}

func openURL(url string) error {
	var cmd string
	var args []string

	url_obj, err := url_tool.Parse(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse url")
	}

	url = url_obj.String()
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		url = strings.ReplaceAll(url, "&", "^&")
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if running under WSL
		if isWSL() {
			// Use 'cmd.exe /c start' to open the URL in the default Windows browser
			cmd = "cmd.exe"
			url = strings.ReplaceAll(url, "&", "^&")
			args = []string{"/c", "start", url}
		} else {
			// Use xdg-open on native Linux environments
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	log.Debug().Msgf("cmd: %s, args: %v", cmd, args)
	return exec.Command(cmd, args...).Start()
}

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
