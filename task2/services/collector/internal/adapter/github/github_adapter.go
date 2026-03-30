package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"repo-stat/services/collector/internal/domain"
)

const requestTimeout = 10 * time.Second

type Adapter struct {
	client *http.Client
}

func NewAdapter() *Adapter {
	return &Adapter{
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (a *Adapter) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Github request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domain.ErrRepoNotFound
	case http.StatusForbidden:
		return nil, domain.ErrGitHubRateLimited
	default:
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: status %d, body: %s", domain.ErrGitHubAPIError, resp.StatusCode, string(body))
	}

	var response struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ForksCount  int32  `json:"forks_count"`
		Stargazers  int32  `json:"stargazers_count"`
		CreatedAt   string `json:"created_at"`
		Visibility  string `json:"visibility"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &domain.Repository{
		Name:        response.Name,
		Description: response.Description,
		Stars:       response.Stargazers,
		Forks:       response.ForksCount,
		CreatedAt:   response.CreatedAt,
		Visibility:  response.Visibility,
	}, nil
}
