package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gabiSmachado/intents/database"
	"github.com/gabiSmachado/intents/datamodel"
	"github.com/gabiSmachado/intents/producer"
	"github.com/gorilla/mux"
)

func main() {
  r := mux.NewRouter()
	r.HandleFunc("/intent", IntentCreate).Methods("POST")
	r.HandleFunc("/intent", IntentList).Methods("GET")
	r.HandleFunc("/intent/{idx}", IntentShow).Methods("GET")
	r.HandleFunc("/intent/{idx}", IntentDelete).Methods("DELETE")
 	 srv := &http.Server{
		Addr:    ":8585",
		Handler: r,
	 } 
	srv.ListenAndServe()	
}

func IntentList(w http.ResponseWriter, r *http.Request) {
	db, _ := database.DBconnect()
	defer db.Close()

	intents, err := database.ListIntents(db)
	if err != nil {
		log.Fatalf("failed to get the intent list")
	}
	json.NewEncoder(w).Encode(intents)
}

func IntentCreate(w http.ResponseWriter, r *http.Request){
	db, _ := database.DBconnect()
	defer db.Close()
	defer r.Body.Close()

	var req datamodel.IntentRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	idx, err := database.Insert(db, req.Intent)
	if err != nil {
		log.Fatalf("failed to insert the intent on datab")
	}

	resp := datamodel.IntentResponse{
		RequestID: req.RequestID,
		IntentID:  idx,
	}
	json.NewEncoder(w).Encode(resp)

	producer.WriteMsg(req.Intent)
}

func IntentShow(w http.ResponseWriter, r *http.Request) {
	db, _ := database.DBconnect()
	defer db.Close()

	idxs := mux.Vars(r)["idx"]
	idx, err := strconv.Atoi(idxs)
	if err != nil {
		log.Fatalf("can't find intent by idx %s", idxs)
	}
	log.Printf("request to view intent %d", idx)
	intent, err := database.IntentShow(db, idx)
	if err != nil {
		log.Fatalf("failed to retrieve intent %d", idx)
	}
	json.NewEncoder(w).Encode(intent)
}

func IntentDelete(w http.ResponseWriter, r *http.Request) {
	db, _ := database.DBconnect()
	defer db.Close()

	idxs := mux.Vars(r)["idx"]
	idx, err := strconv.Atoi(idxs)
	if err != nil {
		log.Fatalf("unable to delete intent %s", idxs)
	}
	log.Printf("deleting intent %d", idx)
	err = database.DeleteIntent(db, idx)
	if err != nil {
		log.Fatalf("unable to delete intent %d", idx)
	}
	w.WriteHeader(200)
}
