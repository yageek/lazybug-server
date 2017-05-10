package bugtracker

import (
	"github.com/andygrunwald/go-jira"
	lazybug "github.com/yageek/lazybug-server/lazybug-protocol"
)

type JiraTracker struct {
	client  *jira.Client
	project string
}

func NewJiraTrackerClient(serverURL, username, password string) (TrackerClient, error) {
	jiraClient, err := jira.NewClient(nil, serverURL)
	if err != nil {
		return nil, err
	}
	jiraClient.Authentication.SetBasicAuth(username, password)
	return &JiraTracker{client: jiraClient}, nil
}

func (j *JiraTracker) CreateTicket(f lazybug.Feedback) error {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: "myuser",
			},
			Reporter: &jira.User{
				Name: "youruser",
			},
			Description: "Test Issue",
			Type: jira.IssueType{
				ID: "60",
			},
			Project: jira.Project{
				Name: "PROJ1",
			},
			Summary: "Just a demo issue",
		},
	}
	_, _, err := j.client.Issue.Create(&i)
	return err
}
