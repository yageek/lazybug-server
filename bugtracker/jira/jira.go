package jira

import (
	"bytes"
	"log"

	"io/ioutil"

	"github.com/andygrunwald/go-jira"
	"github.com/yageek/lazybug-server/bugtracker"
	lazybug "github.com/yageek/lazybug-server/lazybug-protocol"
)

type JiraTracker struct {
	client   *jira.Client
	project  string
	username string
}

func NewJiraTrackerClient(serverURL, username, password, project string) (bugtracker.TrackerClient, error) {
	jiraClient, err := jira.NewClient(nil, serverURL)
	if err != nil {
		return nil, err
	}
	jiraClient.Authentication.SetBasicAuth(username, password)
	return &JiraTracker{client: jiraClient, username: username, project: project}, nil
}

func (j *JiraTracker) CreateTicket(f *lazybug.Feedback) error {

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Reporter: &jira.User{
				Name: j.username,
			},
			Description: f.Content,
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: j.project,
			},
			Summary: "A lazy bug session.",
		},
	}
	log.Printf("Value: %+v \n", i.Fields)
	issue, response, err := j.client.Issue.Create(&i)
	if err != nil {
		log.Printf("Impossible to create ticket:  %q \n", err)
		buff, buffError := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if buffError == nil {
			log.Printf("Response: %+v \n", string(buff))
		}
		return err
	}

	log.Println("SuccessFully created issue:", issue.Key)
	buff := bytes.NewBuffer(f.GetSnapshot())
	_, _, err = j.client.Issue.PostAttachment(issue.ID, buff, f.GetIdentifier()+".jpg")
	if err == nil {
		log.Println("Successfully created attachements")
	}
	return err
}
