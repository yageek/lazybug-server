package store

var (
	feedbackBuckets = []byte("feedbacks")
)

type FeedbackStore interface {
	SaveFeedback(ID string, data []byte) error
	DeleteFeedback(ID string) error
	Close() error
}
