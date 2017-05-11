package bugtracker

import (
	"encoding/json"
	"os"
	"text/template"

	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/yageek/lazybug-server/lazybug-protocol"
)

type templateValues struct {
	Title       string
	Date        string
	Meta        map[string]interface{}
	Description string
}

var templateRaw string = `
{{.Title}} at {{.Date}}
{{.Meta}}
{{.Description}}
`

var tmpl = template.Must(template.New("tmlp").Parse(templateRaw))

type STDOutTracker struct {
	client   *jira.Client
	project  string
	username string
}

func NewSTDOUTClient(serverURL, username, password string) (TrackerClient, error) {
	jiraClient, err := jira.NewClient(nil, serverURL)
	if err != nil {
		return nil, err
	}
	jiraClient.Authentication.SetBasicAuth(username, password)
	return &STDOutTracker{client: jiraClient, username: username}, nil
}

func (s *STDOutTracker) CreateTicket(f *lazybug.Feedback) error {

	var meta map[string]interface{}
	err := json.Unmarshal(f.GetMeta(), &meta)
	if err != nil {
		return err
	}

	values := templateValues{
		Title:       fmt.Sprintf("[LAZYBUG Session] %s", f.GetIdentifier()),
		Date:        f.GetCreationDate(),
		Meta:        meta,
		Description: f.GetContent(),
	}
	return tmpl.Execute(os.Stdout, values)
}
