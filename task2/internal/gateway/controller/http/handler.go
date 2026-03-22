package http

import (
	"encoding/json"
	"net/http"

	"repo-stat/internal/gateway/controller/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client *grpc.Client
}

func NewHandler(client *grpc.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) GetRepoInfo(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	if owner == "" || repo == "" {
		http.Error(w, "owner and repo required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	resp, err := h.client.GetRepoInfo(ctx, owner, repo)
	if err != nil {
		code := http.StatusInternalServerError
		msg := "internal error"

		switch status.Code(err) {
		case codes.NotFound:
			code = http.StatusNotFound
			msg = "repository not found"
		case codes.InvalidArgument:
			code = http.StatusBadRequest
			msg = status.Convert(err).Message()
		case codes.ResourceExhausted:
			code = http.StatusTooManyRequests
			msg = "github rate limit exceeded"
		}

		http.Error(w, msg, code)
		return
	}

	jsonResp := map[string]interface{}{
		"name":             resp.Name,
		"description":      resp.Description,
		"stargazers_count": resp.Stars,
		"forks_count":      resp.Forks,
		"created_at":       resp.CreatedAt,
		"visibility":       resp.Visibility,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResp)
}
