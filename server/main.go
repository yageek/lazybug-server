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
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil || len(buf) < 1 {
			http.Error(w, "Invalid Data", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		buff, err := ioutil.ReadAll(r.Body)
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

		fmt.Printf("Feedb ID: %q \n", feedb)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run()
}

func getDB() *bolt.DB {
	db, err := bolt.Open("lazybugdata.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create feeback buckets
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
