package controller

import (
	"encoding/json"
	"net/http"

	"repo-stat/services/gateway/internal/usecase"
)

type Response struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ForksCount  int32  `json:"forks_count"`
	Stargazers  int32  `json:"stargazers_count"`
	CreatedAt   string `json:"created_at"`
	Visibility  string `json:"visibility"`
}

type Handler struct {
	usecase *usecase.GetRepoInfo
}

func NewHandler(usecase *usecase.GetRepoInfo) *Handler {
	return &Handler{usecase: usecase}
}

// GetRepoInfo godoc
// @Summary      Получить информацию о репозитории
// @Description  Возвращает базовую информацию о публичном репозитории GitHub
// @Tags         repos
// @Accept       json
// @Produce      json
// @Param        owner  path      string  true   "Владелец репозитория (username или организация)"
// @Param        repo   path      string  true   "Название репозитория"
// @Success      200     {object}  map[string]interface{}  "Информация о репозитории"
// @Failure      400     {string}  string                  "Некорректные параметры"
// @Failure      404     {string}  string                  "Репозиторий не найден"
// @Failure      429     {string}  string                  "Превышен лимит запросов к GitHub"
// @Failure      500     {string}  string                  "Внутренняя ошибка"
// @Router       /api/repos/{owner}/{repo} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	if owner == "" || repo == "" {
		http.Error(w, "owner and repo required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	resp, err := h.usecase.Execute(ctx, owner, repo)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response := Response{
		Name:        resp.Name,
		Description: resp.Description,
		ForksCount:  resp.Forks,
		Stargazers:  resp.Stars,
		CreatedAt:   resp.CreatedAt,
		Visibility:  resp.Visibility,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
