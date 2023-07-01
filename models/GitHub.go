package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GitHubRequest represents a GitHub API request model.
type GitHubRequest struct {
	GitHubRepoOwner *string `json:"git_hub_repo_owner,omitempty"`
	GitHubRepoTitle *string `json:"git_hub_repo_title,omitempty"`
	GitHubToken     *string `json:"github_token"`
	Issue           Issue   `json:"issue,omitempty"`
}

// Issue represents a GitHub issue on a repository.
type Issue struct {
	Title     string  `json:"title"`
	Body      *string `json:"body,omitempty"`
	Assignee  *string `json:"assignee,omitempty"`
	Milestone *string `json:"milestone,omitempty"`
	Labels    []Label `json:"labels,omitempty"`
}

// Label represents a GitHub label on an Issue.
type Label struct {
	ID          *int64  `json:"id,omitempty"`
	URL         *string `json:"url,omitempty"`
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
	Default     *bool   `json:"default,omitempty"`
	NodeID      *string `json:"node_id,omitempty"`
}

// Send a request to GitHub to create a new issue.
func (g *GitHubRequest) Send() error {
	requestBody, err := json.Marshal(g.Issue)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", *g.GitHubRepoOwner, *g.GitHubRepoTitle)

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *g.GitHubToken))
	req.Header.Set("x-GitHub-Api-Version", "2022-11-28")

	client := http.Client{Timeout: 10 * time.Second}
	// send the request
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// GenerateIssue returns a newly generated issue.
func (g *GitHubRequest) GenerateIssue(title string, body *string) *Issue {
	return &Issue{
		Title: title,
		Body:  body,
	}
}
