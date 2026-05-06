package github

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

func (a *Adapter) IsRepoExist(ctx context.Context, owner, repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("github request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("github api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return true, nil
}
