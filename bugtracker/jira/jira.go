package bugtracker

import (
	"github.com/andygrunwald/go-jira"
	lazybug "github.com/yageek/lazybug-server/lazybug-protocol"
)

type JiraTracker struct {
	client   *jira.Client
	project  string
	username string
}

func NewJiraTrackerClient(serverURL, username, password string) (TrackerClient, error) {
	jiraClient, err := jira.NewClient(nil, serverURL)
	if err != nil {
		return nil, err
	}
	jiraClient.Authentication.SetBasicAuth(username, password)
	return &JiraTracker{client: jiraClient, username: username}, nil
}

func (j *JiraTracker) CreateTicket(f *lazybug.Feedback) error {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Reporter: &jira.User{
				Name: j.username,
			},
			Description: "Test Issue",
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Name: j.project,
			},
		},
	}
	_, _, err := j.client.Issue.Create(&i)
	return err
}
