package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"github.com/yageek/lazybug-server/bugtracker"
	"github.com/yageek/lazybug-server/lazybug-protocol"
	"github.com/yageek/lazybug-server/store"
	"github.com/yageek/lazybug-server/trackersync"
)

func main() {
	// Init database.
	db, err := store.NewBoltStore("lazybugdata.db")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	//Client
	stdout, err := bugtracker.NewSTDOUTClient("", "", "")
	if err != nil {
		log.Panicln("Err:", err)
	}

	// Sync
	manager := trackersync.NewSyncManager(db, stdout)
	defer manager.Stop()
	manager.Start()

	// Start API
	mux := pat.New()
	mux.Put("/feedbacks", func(w http.ResponseWriter, r *http.Request) {
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil || len(buff) < 1 {
			log.Printf("No data have been provided: %q \n", err)
			http.Error(w, "Invalid Data", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		feedb := &lazybug.Feedback{}

		log.Println("Unmarshalling data...")
		err = proto.Unmarshal(buff, feedb)
		if err != nil {
			log.Printf("Impossible to unmarshal data: %q \n", err)
			http.Error(w, "Invalid Data", http.StatusInternalServerError)
			return
		}

		ID := feedb.GetIdentifier()
		err = db.SaveFeedback(ID, buff)

		if err != nil {
			log.Printf("Impossible to save element: %q \n", err)
			http.Error(w, "Can not save element", http.StatusInternalServerError)
			return
		}
		log.Printf("Feedback %s added. \n", ID)
		w.WriteHeader(http.StatusCreated)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":5555")
}
