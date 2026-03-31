package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"repo-stat/collector/internal/domain"
)

const requestTimeout = 10 * time.Second

type Adapter struct {
	client *http.Client
}

type githubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ForksCount  int32  `json:"forks_count"`
	Stargazers  int32  `json:"stargazers_count"`
	CreatedAt   string `json:"created_at"`
	Visibility  string `json:"visibility"`
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
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("failed to close response body", "error", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domain.ErrRepoNotFound
	case http.StatusForbidden:
		return nil, domain.ErrGitHubRateLimited
	default:
		return nil, fmt.Errorf("%w: status %d, body: %s",
			domain.ErrGitHubAPIError, resp.StatusCode, string(body))
	}

	var gh githubRepo
	if err := json.Unmarshal(body, &gh); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &domain.Repository{
		Name:        gh.Name,
		Description: gh.Description,
		Stars:       gh.Stargazers,
		Forks:       gh.ForksCount,
		CreatedAt:   gh.CreatedAt,
		Visibility:  gh.Visibility,
	}, nil
}
