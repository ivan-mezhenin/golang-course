package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
	"strings"
)

func NewRepoHandler(log *slog.Logger, uc *usecase.GetRepoInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			http.Error(w, `{"error": "url parameter is required"}`, http.StatusBadRequest)
			return
		}

		owner, repo, err := parseGitHubURL(urlParam)
		if err != nil {
			http.Error(w, `{"error": "invalid github url format"}`, http.StatusBadRequest)
			return
		}

		repoInfo, err := uc.Execute(r.Context(), owner, repo)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
			return
		}

		response := dto.RepoResponse{
			Name:        repoInfo.Name,
			Description: repoInfo.Description,
			Stars:       repoInfo.Stars,
			Forks:       repoInfo.Forks,
			CreatedAt:   repoInfo.CreatedAt,
			Visibility:  repoInfo.Visibility,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode repo response", "error", err)
		}
	}
}

func parseGitHubURL(rawURL string) (owner, repo string, err error) {
	clean := strings.TrimPrefix(rawURL, "https://github.com/")
	clean = strings.TrimPrefix(clean, "http://github.com/")
	clean = strings.Trim(clean, "/")

	parts := strings.Split(clean, "/")
	if len(parts) < 2 {
		return "", "", domain.ErrInvalidInput
	}

	return parts[0], parts[1], nil
}
