package github

import (
	"context"

	"strings"

	"github.com/google/go-github/github"
	"github.com/yageek/lazybug-server/bugtracker"
	lazybug "github.com/yageek/lazybug-server/lazybug-protocol"
	"golang.org/x/oauth2"
)

type GithubTracker struct {
	client *github.Client
	repo   string
	user   string
}

func NewGithubTracker(user, repo, token string) (bugtracker.TrackerClient, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	_, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		return nil, err
	}
	return &GithubTracker{client: client, repo: repo, user: user}, nil
}

func (g *GithubTracker) CreateTicket(f *lazybug.Feedback) error {
	ctx := context.Background()

	title := strings.SplitAfterN(f.GetContent(), "", 10)[0] + "..."
	request := &github.IssueRequest{
		Title: &title,
	}
	_, _, err := g.client.Issues.Create(ctx, g.user, g.repo, request)
	return err
}
