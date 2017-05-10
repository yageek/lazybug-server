package store

import (
	"time"

	"github.com/boltdb/bolt"
)

type boltStore struct {
	db *bolt.DB
}

func NewBoltStore(dir string) (FeedbackStore, error) {
	db, err := bolt.Open("lazybugdata.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return &boltStore{db: db}, nil
}

func (s *boltStore) SaveFeedback(ID string, data []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(feedbackBuckets)
		if err != nil {
			return err
		}
		return b.Put([]byte(ID), data)
	})
	return err
}

func (s *boltStore) DeleteFeedback(ID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(feedbackBuckets))
		return b.Delete([]byte(ID))
	})
}
func (s *boltStore) Close() error {
	return s.db.Close()
}
