package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Repository struct {
	Name        string    `json:"name"`
	Visibility  string    `json:"visibility"`
	Description string    `json:"description"`
	Forks       int       `jsom:"forks_count"`
	Stars       int       `json:"stargazers_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	var repoURL string

	flag.StringVar(&repoURL, "url", "", "URL of repo (example: https://github.com/owner/repo)")
	flag.Parse()

	if repoURL == "" {
		if len(os.Args) > 1 {
			repoURL = os.Args[1]
		}
	}

	if repoURL == "" {
		fmt.Println("Error: Enter URL of repo")
		os.Exit(1)
	}

	owner, repo, err := parseGitHubRepoURL(repoURL)
	if err != nil {
		fmt.Printf("Incorrect URL: %v\n", err)
		os.Exit(1)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while getting repository: ", err)
		return
	}

	switch response.StatusCode {
	case http.StatusOK:
		var repo Repository
		if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
			fmt.Println("Error while parsing json: ", err)
			return
		}

		fmt.Printf("\033[35m"+"Name: %s\n", repo.Name)
		fmt.Printf("Visibility: %s\n", repo.Visibility)
		fmt.Printf("Description: %s\n", repo.Description)
		fmt.Printf("Amount of forks: %d\n", repo.Forks)
		fmt.Printf("Amount of stars: %d\n", repo.Stars)
		fmt.Printf("Created at: %s\n", repo.CreatedAt.Format("2006-01-02 15:04:05 MST"))

	case http.StatusNotFound:
		fmt.Println("Repository not found")

	case http.StatusMovedPermanently:
		fmt.Println("Repository was moved")
		location := response.Header.Get("Location")
		if location != "" {
			fmt.Printf("New location: %s\n", location)
		}
	default:
		body, _ := io.ReadAll(response.Body)
		fmt.Printf("Unexpected status %d for %s/%s\n%s\n", response.StatusCode, owner, repo, string(body))
	}
}

func parseGitHubRepoURL(s string) (owner, repo string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", "", err
	}

	if u.Host != "github.com" && u.Host != "www.github.com" {
		return "", "", fmt.Errorf("Expected github.com in URL")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("Expected /owner/repo in URL")
	}

	owner = parts[0]
	repo = parts[1]

	repo = strings.TrimSuffix(repo, ".git")

	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("Could't define owner or repo")
	}

	return owner, repo, nil
}
