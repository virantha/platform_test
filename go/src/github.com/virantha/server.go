package main

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
)

type MessageJSON struct {
    Message string `json:"message"`
    Recipients []string `json:recipients"`
}

type RouteJSON struct {
    IP string `json:"ip"`
    Recipients []string `json:"recipients"`
}

type RoutesJSON struct {
    Message string `json:"message"`
    Routes []RouteJSON `json:"routes"` 
}

func divide(msg_length int, categories []int) (msgs map[int]int){
    var remainder int
    msgs = make(map[int]int) // Hash of category(msgs/request) to number of such routes

    // Divide the message count into the bins
    // starting with the largest bin first (25 msgs/req)
    // Assume bins are in descending order
    // Assume always have a bin of 1! 
    remainder = msg_length

    i:=0

    sum := 0  // keep summing up as a sanity check
    for remainder > 0 {
        messages := remainder/categories[i]
        remainder = remainder % categories[i]
        msgs[categories[i]] = messages
        sum += categories[i]*messages
        i = i+1
    }
    fmt.Println(msgs)
    fmt.Printf("Sum: %d, start:%d", sum, msg_length)
    return msgs
}

func GetMessageAllocation(message string, recipients []string) (error, *RoutesJSON){
    allocation := divide(len(recipients), []int{25,10,5,1})
    // Now, take the allocation and split up the recipients into multiple slices
    IPs := map[int]string{
        25: "10.0.4.",
        10: "10.0.3.",
        5 : "10.0.2.",
        1 : "10.0.1.",
    }
    var routesJSON RoutesJSON
    routesJSON.Message = message

    i := 0
    for bin,bin_count := range(allocation) {
        // For each bin, slice up the recipients
        // i will keep track of where we are in the recipients
        for bin_index:=0; bin_index<bin_count; bin_index++ {
            routeJSON := RouteJSON{ IP: fmt.Sprintf("%s%d", IPs[bin], bin_index+1),
                                    Recipients: recipients[i:i+bin],
                                  }
            routesJSON.Routes = append(routesJSON.Routes, routeJSON)
            i = i+bin
        }
    }
    return nil, &routesJSON
}



func MessageRouter(w http.ResponseWriter, r *http.Request) {
    var messageJSON MessageJSON

    log.Println("Got new message")
    err := json.NewDecoder(r.Body).Decode(&messageJSON)
    if err != nil || messageJSON.Message == "" || messageJSON.Recipients == nil{
        WriteJSONError(w, 400, "Malformed request")
        return
    }
    // Extract the POST message components
    message := messageJSON.Message
    recipients := messageJSON.Recipients

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

func main() {

    // Setup the routes
    r := mux.NewRouter()
    r.HandleFunc("/message/route", MessageRouter).Methods("POST")
    http.Handle("/message/", r)

    // Start the server
    log.Println("Starting testing")
    http.ListenAndServe(":8080", nil)


}
