package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy/lfs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type GoogleDriveLFSHandler struct {
	driveClient *GoogleDriveClient
	ctx         context.Context
	proxyName   string
}

func NewGoogleDriveLFSHandler(client *GoogleDriveClient, ctx context.Context, proxyName string) lfs.BasicLFSHandler {
	return &GoogleDriveLFSHandler{
		driveClient: client,
		ctx:         ctx,
		proxyName:   proxyName,
	}
}

func (h *GoogleDriveLFSHandler) HandleBatchAPI(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	decoder := json.NewDecoder(r.Body)
	batchBody := &lfs.LFSBatchRequest{}
	err := decoder.Decode(batchBody)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding batch request")
		http.Error(w, "Error decoding batch request", http.StatusBadRequest)
		return
	}

	numObjects := len(batchBody.Objects)
	objects := make([]any, numObjects)

	switch batchBody.Operation {
	case "download":
		for i, obj := range batchBody.Objects {
			// download object
			objects[i] = &lfs.LFSBatchResponseObject{
				Oid:  obj.Oid,
				Size: obj.Size,
				Actions: map[string]lfs.LFSBatchResponseAction{
					"download": {
						Href:      h.replaceUserRepo(fmt.Sprintf("%s%s/objects/%s", config.Global.BaseApiUrl, h.getApiUrlBase(), obj.Oid), user, repo),
						ExpiresAt: time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
					},
				},
				HashAlgorithm: batchBody.HashAlgorithm,
			}
		}
	case "upload":
		for i, obj := range batchBody.Objects {
			// upload object
			objects[i] = &lfs.LFSBatchResponseObject{
				Oid:  obj.Oid,
				Size: obj.Size,
				Actions: map[string]lfs.LFSBatchResponseAction{
					"upload": {
						Href:      h.replaceUserRepo(fmt.Sprintf("%s%s/objects/%s", config.Global.BaseApiUrl, h.getApiUrlBase(), obj.Oid), user, repo),
						ExpiresAt: time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
					},
				},
				HashAlgorithm: batchBody.HashAlgorithm,
			}
		}
	default:
		log.Error().Msg("Invalid operation")
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return
	}

	ret := &lfs.LFSBatchResponse{
		Transfer: "basic",
		Objects:  objects,
	}

	log.Debug().Interface("batch response", ret).Msg("Batch response")

	w.Header().Set("Content-Type", lfs.LFS_HEADER)
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Error().Err(err).Msg("Error encoding batch response")
		http.Error(w, "Error encoding batch response", http.StatusInternalServerError)
		return
	}

}

func (h *GoogleDriveLFSHandler) HandleDownloadAPI(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("HandleDownloadAPI")
}

func (h *GoogleDriveLFSHandler) HandleUploadAPI(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("HandleUploadAPI")
}

func (h *GoogleDriveLFSHandler) HandleVerifyAPI(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("HandleVerifyAPI")
}

func (h *GoogleDriveLFSHandler) HandleLockAPI(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("HandleLockAPI")
}

func (h *GoogleDriveLFSHandler) HandleUnlockAPI(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("HandleUnlockAPI")
}

func (h *GoogleDriveLFSHandler) HandleVerifyLockAPI(w http.ResponseWriter, r *http.Request) {
	ret := make(map[string]any)
	ret["ours"] = []any{}
	ret["theirs"] = []any{}

	w.Header().Set("Content-Type", lfs.LFS_HEADER)
	err := json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Error().Err(err).Msg("Error encoding verify lock response")
		http.Error(w, "Error encoding verify lock response", http.StatusInternalServerError)
		return
	}
}

func (h *GoogleDriveLFSHandler) getApiUrlBase() string {
	return fmt.Sprintf("/%s/{user}/{repo}/lfs", h.proxyName)
}

func (h *GoogleDriveLFSHandler) replaceUserRepo(template string, user string, repo string) string {
	replacer := strings.NewReplacer("{user}", user, "{repo}", repo)
	return replacer.Replace(template)
}
