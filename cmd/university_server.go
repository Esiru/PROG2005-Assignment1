package main

import (
	"RESTUniversity"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"encoding/json"
)

//a global variabel that stores when the service started, i.e. when the executive was launched
var startTime time.Time

/*
	The handler for the /diag/ endpoint. 
	Delivers the http-status of the two APIs queried by unisearcher,
	in addition to unisearcher's uptime and version
	Sends two http-requests, one to http://universities.hipolabs.com/,
	and one to ttps://restcountries.com/, getting a relatively small package from both
	This to ascertain whether the services are running or not
*/

func handlerDiag(w http.ResponseWriter, r *http.Request) {
	//A diag-struct that encapsulates the information the endpoint will provide
	//unistatus, countrystatus, version, uptime
	var s RESTUniversity.Diag
	//The first http-request, this to http://universities.hipolabs.com/
	//hoping to find all universities with norwegian in the name
	//see Client func in client.go
	res, err := RESTUniversity.Client(RESTUniversity.UNIVERSITIES_PATH+"norwegian")
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
		return
	}
	//the diag unistatus is filled with the status-code from hipolabs 
	s.UniStatus = res.StatusCode
	//The second http-reqeust, this to restcountries
	//tries to get only the names from all universities in Norway
	res, err = RESTUniversity.Client("https://restcountries.com/v3.1/name/norway?fields=name")
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
		return
	}
	//the diag-struct is further initialized with countrystatus, version and uptime
	s.CountryStatus = res.StatusCode
	s.Version = "v1"
	//Finds how much time has passed since startTime was filled, and converts that time.Duration-type
	//to float64
	s.Uptime = time.Duration(time.Since(startTime)).Seconds()

	http.Header.Add(w.Header(), "content-type", "application/json")
	
	json.NewEncoder(w).Encode(s)
}

func main() {
	//Fills starttime with current date 
	startTime = time.Now()
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has been set. Default: " + RESTUniversity.DEFAULT_PORT)
		port = RESTUniversity.DEFAULT_PORT
	}

	//The various handlers
	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"diag", handlerDiag)
	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"diag"+"/", handlerDiag)
	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"uniinfo/", RESTUniversity.HandlerUniversity())
	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"neighbourunis/", RESTUniversity.HandlerUniversity())

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}