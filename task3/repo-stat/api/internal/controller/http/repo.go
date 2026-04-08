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

// NewRepoHandler godoc
// @Summary      Получить информацию о репозитории по URL
// @Description  Принимает URL GitHub репозитория и возвращает его основные данные
// @Tags         repositories
// @Accept       json
// @Produce      json
// @Param        url  query  string  true  "URL репозитория (например: https://github.com/octocat/Hello-World)"
// @Success      200  {object}  dto.RepoResponse
// @Failure      400  {object}  map[string]string  "Некорректный URL или параметры"
// @Failure      404  {object}  map[string]string  "Репозиторий не найден"
// @Failure      429  {object}  map[string]string  "Превышен лимит запросов GitHub"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /api/repositories/info [get]
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
			switch err {
			case domain.ErrInvalidInput:
				http.Error(w, `{"error": "invalid owner or repo name"}`, http.StatusBadRequest)
			case domain.ErrRepoNotFound:
				http.Error(w, `{"error": "repository not found"}`, http.StatusNotFound)
			case domain.ErrGitHubRateLimited:
				http.Error(w, `{"error": "github rate limit exceeded"}`, http.StatusTooManyRequests)
			default:
				log.Error("failed to get repository info", "error", err, "owner", owner, "repo", repo)
				http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			}
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
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
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
