package domain

type RepoRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type RepoResponse struct {
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
	Visibility  string `json:"visibility"`
	Error       string `json:"error,omitempty"`
}
