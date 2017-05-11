package trackersync

import (
	"log"
	"time"

	"github.com/yageek/lazybug-server/bugtracker"
	lazybug "github.com/yageek/lazybug-server/lazybug-protocol"
	"github.com/yageek/lazybug-server/store"
)

type SyncManager struct {
	ticker  *time.Ticker
	store   store.FeedbackStore
	tracker bugtracker.TrackerClient
}

func NewSyncManager(st store.FeedbackStore, client bugtracker.TrackerClient) *SyncManager {
	t := time.NewTicker(30 * time.Second)
	return &SyncManager{ticker: t, store: st, tracker: client}
}

func (s *SyncManager) Start() {
	go func() {
		for range s.ticker.C {
			s.performSync()
		}
	}()
}

func (s *SyncManager) Stop() {
	s.ticker.Stop()
}

func (s *SyncManager) performSync() {
	log.Println("Start sync pass...")

	var elementsToDelete []string
	s.store.Iterate(store.IteratorFunc(func(err error, feedback *lazybug.Feedback) {
		if err != nil {
			log.Printf("Impossible to fetch feedback: %q \n", err)
			return
		}
		err = s.tracker.CreateTicket(feedback)
		if err == nil {
			elementsToDelete = append(elementsToDelete, feedback.GetIdentifier())
		} else {
			log.Printf("Error creating ticket: %q \n", err)
		}
	}))

	if len(elementsToDelete) > 0 {
		log.Printf("Elements to Delete: %d \n", len(elementsToDelete))
		if err := s.store.DeleteFeedbacks(elementsToDelete); err != nil {
			log.Printf("Impossible to delete elements: %q \n", err)
		}
	}

}
