package store

import lazybug "github.com/yageek/lazybug-server/lazybug-protocol"

var (
	feedbackBuckets = []byte("feedbacks")
)

type Iterator interface {
	Next(error, *lazybug.Feedback)
}
type IteratorFunc func(error, *lazybug.Feedback)

func (f IteratorFunc) Next(err error, feedb *lazybug.Feedback) {
	f(err, feedb)
}

type Cursor interface {
	Next() lazybug.Feedback
}
type FeedbackStore interface {
	SaveFeedback(ID string, data []byte) error
	DeleteFeedbacks(IDs []string) error
	Iterate(i Iterator)
	Close() error
}
