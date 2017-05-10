package bugtracker

import "github.com/google/go-github/github"

type GithubTracker struct {
	client  *github.Client
	project string
	user    string
}

func NewGithubTracker(token string) (TrackerClient, error) {
	return nil, nil
}
