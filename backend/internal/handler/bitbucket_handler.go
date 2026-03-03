package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/igor-trentini/diffable/backend/internal/bitbucket"
	"github.com/igor-trentini/diffable/backend/internal/cache"
)

type BitbucketHandler struct {
	bbClient bitbucket.Client
	cache    cache.Cache
}

func NewBitbucketHandler(bbClient bitbucket.Client, c cache.Cache) *BitbucketHandler {
	return &BitbucketHandler{
		bbClient: bbClient,
		cache:    c,
	}
}

type repoResponse struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func (h *BitbucketHandler) ListRepositories(w http.ResponseWriter, r *http.Request) {
	workspace := r.URL.Query().Get("workspace")
	if workspace == "" {
		writeErrorJSON(w, http.StatusBadRequest, "validation_error", "workspace query parameter is required")
		return
	}

	query := r.URL.Query().Get("q")
	cacheKey := "repos:" + workspace

	// Try cache for the full (unfiltered) list
	var repos []bitbucket.Repository
	if cached, ok := h.cache.Get(cacheKey); ok {
		var cachedRepos []bitbucket.Repository
		if err := json.Unmarshal([]byte(cached), &cachedRepos); err == nil {
			repos = cachedRepos
		}
	}

	if repos == nil {
		var err error
		repos, err = h.bbClient.ListRepositories(r.Context(), workspace)
		if err != nil {
			slog.Error("failed to list repositories", "workspace", workspace, "error", err)
			writeErrorJSON(w, http.StatusBadGateway, "external_service_error", "Failed to list repositories from Bitbucket")
			return
		}

		// Cache the full list
		if data, err := json.Marshal(repos); err == nil {
			h.cache.Set(cacheKey, string(data), 5*time.Minute)
		}
	}

	// Apply query filter
	if query != "" {
		query = strings.ToLower(query)
		var filtered []bitbucket.Repository
		for _, repo := range repos {
			if strings.Contains(strings.ToLower(repo.Slug), query) ||
				strings.Contains(strings.ToLower(repo.Name), query) {
				filtered = append(filtered, repo)
			}
		}
		repos = filtered
	}

	response := make([]repoResponse, 0, len(repos))
	for _, repo := range repos {
		response = append(response, repoResponse{
			Slug:     repo.Slug,
			Name:     repo.Name,
			FullName: repo.FullName,
		})
	}

	writeJSON(w, http.StatusOK, response)
}
