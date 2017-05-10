package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"github.com/yageek/lazybug-server/lazybug-protocol"
)

var (
	feedbackBuckets = []byte("feedbacks")
)

func main() {
	db := getDB()
	defer db.Close()

	mux := pat.New()
	mux.Put("/feedbacks", func(w http.ResponseWriter, r *http.Request) {
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil || len(buff) < 1 {
			http.Error(w, "Invalid Data", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		fmt.Printf("Output Length: %d \n", len(buff))
		if err != nil {
			http.Error(w, "Invalid Data", http.StatusBadRequest)
			return
		}

		feedb := &lazybug.FeedbackAddRequest{}

		err = proto.Unmarshal(buff, feedb)
		if err != nil {
			http.Error(w, "Invalid Data", http.StatusInternalServerError)
			return
		}
		err = db.Update(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(feedbackBuckets))
			return b.Put([]byte(feedb.GetIdentifier()), buff)
		})

		if err != nil {
			http.Error(w, "Can not save element", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Feedb ID %q added!\n", feedb.GetContent())
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":5555")
}

func getDB() *bolt.DB {
	db, err := bolt.Open("lazybugdata.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create feedback buckets
	tx, err := db.Begin(true)
	if err != nil {
		panic(err)
	}

	_, err = tx.CreateBucketIfNotExists(feedbackBuckets)
	if err != nil {
		panic(err)
	}

	return db
}
