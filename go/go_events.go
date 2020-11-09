package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var daprPort = os.Getenv("DAPR_HTTP_PORT")

const stateStoreName = `events`

var stateURL = fmt.Sprintf(`http://localhost:%s/v1.0/state/%s`, daprPort, stateStoreName)

// Event represents an event, be it meetings, birthdays etc
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

	resp, err := http.Post(stateURL, "application/json", bytes.NewBuffer(state))
	if err != nil {
		log.Fatalln("Error posting to state", err)
		http.Error(w, "Failed to write to store", http.StatusServiceUnavailable)
	}
	log.Printf("Response after posting to state: %s", string(resp.Status))
	http.Error(w, "All Okay", http.StatusOK)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	var id string
	json.NewDecoder(r.Body).Decode(&id)

	deleteURL := stateURL + "/" + id

	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error deleting event", err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	log.Printf(string(bodyBytes))
}
