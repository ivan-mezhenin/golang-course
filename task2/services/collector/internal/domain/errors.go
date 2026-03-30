package domain

import "errors"

var (
	ErrRepoNotFound      = errors.New("Repository not found")
	ErrInvalidInput      = errors.New("Invalid owner or repository name")
	ErrGitHubAPIError    = errors.New("Github api returned an error")
	ErrGitHubRateLimited = errors.New("Github rate limit exceeded")
)
