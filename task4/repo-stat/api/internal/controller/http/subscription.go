package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
	"strings"
)

func NewCreateSubscriptionHandler(log *slog.Logger, uc *usecase.SubscriptionUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.SubscriptionResponse

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", "error", err)
			http.Error(w, `{"error": "invalid json body"}`, http.StatusBadRequest)
			return
		}

		if req.Owner == "" || req.Repo == "" {
			http.Error(w, `{"error": "owner and repo are required"}`, http.StatusBadRequest)
			return
		}

		err := uc.Create(r.Context(), req.Owner, req.Repo)
		if err != nil {
			log.Error("failed to create subscription", "error", err, "owner", req.Owner, "repo", req.Repo)

			if strings.Contains(err.Error(), "already exists") {
				http.Error(w, `{"error": "subscription already exists"}`, http.StatusConflict)
				return
			}

			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "subscription created successfully",
		}); err != nil {
			log.Error("failed to write subscription creation success", "error", err)
		}
	}
}

func NewDeleteSubscriptionHandler(log *slog.Logger, uc *usecase.SubscriptionUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := strings.TrimSpace(r.PathValue("owner"))
		repo := strings.TrimSpace(r.PathValue("repo"))

		if owner == "" || repo == "" {
			http.Error(w, `{"error": "owner and repo are required"}`, http.StatusBadRequest)
			return
		}

		err := uc.Delete(r.Context(), owner, repo)
		if err != nil {
			log.Error("failed to delete subscription", "error", err, "owner", owner, "repo", repo)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "subscription deleted successfully",
		}); err != nil {
			log.Error("failed to write removal success", "error", err)
		}
	}
}

func NewListSubscriptionsHandler(log *slog.Logger, uc *usecase.SubscriptionUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := uc.List(r.Context())
		if err != nil {
			log.Error("failed to list subscriptions", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		subscriptions := dto.ListSubscriptionsResponse{Subscriptions: make([]dto.SubscriptionResponse, 0, len(resp))}
		for _, sub := range resp {
			subscriptions.Subscriptions = append(subscriptions.Subscriptions, dto.SubscriptionResponse{
				Owner: sub.Owner,
				Repo:  sub.Repo,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
			log.Error("failed to write list subscriptions response", "error", err)
		}
	}
}
