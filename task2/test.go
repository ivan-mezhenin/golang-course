package main

import (
	"net/http"
	"time"
)

type Response struct {
	StatusCode  int       `json:"status_code"`
	Name        string    `json:"name"`
	Visibility  string    `json:"visibility"`
	Description string    `json:"description"`
	Forks       int       `jsom:"forks_count"`
	Stars       int       `json:"stargazers_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/repos/{owner}/{repo}", handleGetRepositoryInfo)

}

func handleGetRepositoryInfo(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	if owner == "" || repo == "" {
		http.Error(w, "owner and required", http.StatusBadRequest)
	}

}
