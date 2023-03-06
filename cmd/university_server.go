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
var startTime time.Time

func handlerDiag(w http.ResponseWriter, r *http.Request) {
	var s RESTUniversity.Diag
	res, err := RESTUniversity.Client(RESTUniversity.UNIVERSITIES_PATH+"norwegian")
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
		return
	}
	s.UniStatus = res.StatusCode
	s.Version = "v1"
	res, err = RESTUniversity.Client("https://restcountries.com/v3.1/name/norway?fields=name")
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
		return
	}
	s.CountryStatus = res.StatusCode
	s.Uptime = time.Duration(time.Since(startTime)).Seconds()

	http.Header.Add(w.Header(), "content-type", "application/json")
	//s := RESTUniversity.Diag(res1.Status, res2.StatusCode, "v1", time.Since(startTime))
	
	json.NewEncoder(w).Encode(s)
}

func main() {
	startTime = time.Now()
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has been set. Default: " + RESTUniversity.DEFAULT_PORT)
		port = RESTUniversity.DEFAULT_PORT
	}

	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"diag", handlerDiag)
	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"diag"+"/", handlerDiag)
	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"uniinfo/", RESTUniversity.HandlerUniversity())
	http.HandleFunc(RESTUniversity.UNISEARCHER_PATH+"neighbourunis/", RESTUniversity.HandlerUniversity())

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}