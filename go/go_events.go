package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var daprPort = os.Getenv("DAPR_HTTP_PORT")

const stateStoreName = `events`

var stateUrl = fmt.Sprintf(`http://localhost:%s/v1.0/state/%s`, daprPort, stateStoreName)

type Event struct {
	name string
	date string
	id   string
}

func addEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	json.NewDecoder(r.Body).Decode(&event)
	log.Printf("Event Name: %s", event.name)
	log.Printf("Event Date: %s", event.date)
	log.Printf("Event ID: %s", event.id)

	state, _ := json.Marshal(map[string]string{
		"key":   event.id,
		"value": event.name + " " + event.date,
	})

	resp, err := http.Post(stateUrl, "application/json", bytes.NewBuffer(state))
	if err != nil {
		log.Fatalln("Error posting to state", err)
		http.Error(w, "Failed to write to store", http.StatusServiceUnavailable)
	}
	log.Printf("Response after posting to state: %s", resp)
	http.Error(w, "All Okay", http.StatusOK)
}
