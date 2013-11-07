package main

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"errors"
	"regexp"
)

// Some globals that really should be in a configuration file or member data
const MAX_RECIPIENTS = 5000
// Our table of routes
var IPs map[int]string = map[int]string{
	25: "10.0.4.",
	10: "10.0.3.",
	5:  "10.0.2.",
	1:  "10.0.1.",
}
// Calculate the reverse sorted keys of the route categories
var Categories []int = reverse_sort_keys(IPs)
// Compiled regex to match a 10-digit phone number
var PhoneRe *regexp.Regexp = regexp.MustCompile(`\d{10}`)

// The incoming request
type MessageJSON struct {
	Message    string   `json:"message"`
	Recipients []string `json:recipients"`
}

// A single route for the response
type RouteJSON struct {
	IP         string   `json:"ip"`
	Recipients []string `json:"recipients"`
}

// The complete response with the message and array of routes
type RoutesJSON struct {
	Message string      `json:"message"`
	Routes  []RouteJSON `json:"routes"`
}


// This function divides a number into the given throughput categories.
// It returns a hash mapping each throughput category to the number of
// such categories required to satisfy sendint the total number in the
// least number of responses
func divide(msg_length int, categories []int) (msgs map[int]int) {
	var remainder int
	msgs = make(map[int]int) // Hash of category(msgs/request) to number of such routes
	for _, k := range categories {
		msgs[k] = 0 // Make sure we initialize each category to 0
	}

	// Divide the message count into the bins
	// starting with the largest bin first (25 msgs/req)
	// Assume bins are in descending order
	// Assume always have a bin of 1! 
	remainder = msg_length

	i := 0

	sum := 0 // keep summing up as a sanity check
	for remainder > 0 {
		messages := remainder / categories[i]
		remainder = remainder % categories[i]
		msgs[categories[i]] = messages
		sum += categories[i] * messages
		i = i + 1
	}
	return msgs
}

// At some point, I really need to write a sorted dict class
func reverse_sort_keys(categories map[int]string) (keys []int) {
	for key, _ := range categories {
		keys = append(keys, key)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	return keys
}

// Create the JSON return object with the routes
func GetMessageAllocation(message string, recipients []string) (error, *RoutesJSON) {

	if len(recipients) > MAX_RECIPIENTS {
		err := errors.New(fmt.Sprintf("Got %d recipients, but maximum allowed is %d", len(recipients), MAX_RECIPIENTS))
		return err, nil
	}

	for _, recipient := range recipients {
		// Check that recipient is a valid phone number
		if !PhoneRe.MatchString(recipient) {
			err := errors.New(fmt.Sprintf("Invalid phone number %s", recipient))
			return err, nil
		}
	}

	// Now, take the allocation and split up the recipients into multiple slices
	// First, we need to get a reverse sorted list of the categories, so that we allocate
	// the recipients in the most efficient (least number of requests) manner.
	allocation := divide(len(recipients), Categories)

	var routesJSON RoutesJSON // The return object
	routesJSON.Message = message

	i := 0
	for bin, bin_count := range allocation {
		// For each bin, slice up the recipients
		// i will keep track of where we are in the recipients
		for bin_index := 0; bin_index < bin_count; bin_index++ {
			routeJSON := RouteJSON{
				IP:         fmt.Sprintf("%s%d", IPs[bin], bin_index+1),
				Recipients: recipients[i : i+bin],
			}
			routesJSON.Routes = append(routesJSON.Routes, routeJSON)
			i = i + bin
		}
	}
	return nil, &routesJSON
}

// The message/route POST endpoint handler.  This function
// decodes the incoming message, calls the allocation routine,
// and sets the JSON response.
func MessageRouter(w http.ResponseWriter, r *http.Request) {
	var messageJSON MessageJSON

	log.Println("Got new message")
	err := json.NewDecoder(r.Body).Decode(&messageJSON)
	switch { // Do some sanity checking on the JSON we received
	case err != nil:
		WriteJSONError(w, 400, "Malformed request")
		return
	case messageJSON.Message == "":
		WriteJSONError(w, 400, "Message cannot be empty")
		return
	case messageJSON.Recipients == nil:
		WriteJSONError(w, 400, "Recipients list cannot be empty")
		return
	}

	// Extract the POST message components
	message := messageJSON.Message
	recipients := messageJSON.Recipients

	// TODO: check if unique phone numbers
	// Call the allocation routine that returns a RoutesJSON object
	var routesJSON *RoutesJSON
	err, routesJSON = GetMessageAllocation(message, recipients)
	if err != nil {
		WriteJSONError(w, 400, err.Error())
		log.Printf("ERROR: 400: %s", err.Error())
		return
	}

	// Now, create the JSON response
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(routesJSON)
	if err != nil {
		WriteJSONError(w, 500, "Server malfunction")
		return
	} else {
		w.Write(j)
		log.Println("Provided json")
	}
}

// Main entry point that sets a single endpoint and starts the
// server
func main() {

	// Setup the routes
	r := mux.NewRouter()
	r.HandleFunc("/message/route", MessageRouter).Methods("POST")
	http.Handle("/message/", r)

	// Start the server
	log.Println("Starting testing")
	http.ListenAndServe(":8080", nil)

}
