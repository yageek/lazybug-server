package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/yageek/lazybug-server/bugtracker"
	lazybug "github.com/yageek/lazybug-server/lazybug-protocol"
)

type templateValues struct {
	Title       string
	Date        string
	Meta        string
	Description string
}

var templateRaw string = `
{{.Title}} at {{.Date}}
{panel:title=Meta}
{{.Meta}}
{panel}
{{.Description}}
`

var tmpl = template.Must(template.New("tmlp").Parse(templateRaw))

type JiraTracker struct {
	client   *jira.Client
	project  string
	username string
}

func NewTrackerClient(serverURL, username, password, project string) (bugtracker.TrackerClient, error) {
	jiraClient, err := jira.NewClient(nil, serverURL)
	if err != nil {
		return nil, err
	}
	jiraClient.Authentication.SetBasicAuth(username, password)
	return &JiraTracker{client: jiraClient, username: username, project: project}, nil
}

func (j *JiraTracker) CreateTicket(f *lazybug.Feedback) error {

	var meta map[string]interface{}
	err := json.Unmarshal(f.GetMeta(), &meta)
	if err != nil {
		return err
	}

	content := strings.Replace(FormatMeta(meta), "\n\n", "\n", -1)
	values := templateValues{
		Title:       fmt.Sprintf("[LAZYBUG Session] %s", f.GetIdentifier()),
		Date:        f.GetCreationDate(),
		Meta:        content,
		Description: f.GetContent(),
	}

	contentbuff := bytes.NewBufferString("")
	err = tmpl.Execute(contentbuff, values)
	if err != nil {
		return err
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Reporter: &jira.User{
				Name: j.username,
			},
			Description: contentbuff.String(),
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: j.project,
			},
			Summary: "A lazy bug session.",
		},
	}

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
