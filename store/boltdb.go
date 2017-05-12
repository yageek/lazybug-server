package store

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/micro/protobuf/proto"
	lazybug "github.com/yageek/lazybug-server/lazybug-protocol"
)

type boltStore struct {
	db *bolt.DB
}

func NewBoltStore(dir string) (FeedbackStore, error) {
	db, err := bolt.Open("lazybugdata.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists([]byte(feedbackBuckets))
	if err != nil {
		return nil, err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &boltStore{db: db}, nil
}

func (s *boltStore) SaveFeedback(ID string, data []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(feedbackBuckets))
		return b.Put([]byte(ID), data)
	})
	return err
}

func (s *boltStore) DeleteFeedbacks(IDs []string) error {
	if len(IDs) < 1 {
		return nil
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(feedbackBuckets))

		for _, v := range IDs {

			if err := b.Delete([]byte(v)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *boltStore) Close() error {
	return s.db.Close()
}

func (s *boltStore) Iterate(i Iterator) {
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(feedbackBuckets))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			feedb := &lazybug.Feedback{}
			err := proto.Unmarshal(v, feedb)
			i.Next(err, feedb)
		}
		return nil
	})
}
