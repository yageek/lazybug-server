package bugtracker

import lazybug "github.com/yageek/lazybug-server/lazybug-protocol"

type TrackerClient interface {
	CreateTicket(f *lazybug.Feedback) error
}
