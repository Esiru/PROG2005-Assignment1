package RESTUniversity

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var countries []Country

var universities []University

//The first handler. Responsible for ascertaining whether a function is get, or not
//if yes, then the program proceeds, if not then it returns and another request must be sent
func HandlerUniversity() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleUniversityGet(w, r)
		default:
			http.Error(w, "Method not allowed, only GET requests are supported", http.StatusNotImplemented)
		}
	}
}

//If the request is get, this function handles wheter it queries all universites that match
//the queried name (partially or completely), 
//or all matching universites in a given country's bordering countries
func handleUniversityGet(w http.ResponseWriter, r *http.Request) {
	//If two slices, one containing countries and the other unis,
	//have not been filled, then the function gets ALL information from
	//the two attendant APIs
	//While this might represent a large start-up cost,
	//it is MUCH faster than sending several requests per country or university
	if len(countries) == 0 {
		res, err := Client(COUNTRIES_ALL_PATH)
		if err != nil {
			return
		}
		arrGen(res, &countries)
	}
	//Both of the slices are filled with Client and arrGen-funcs from client.go
	if len(universities) == 0 {
		res, err := Client(UNIVERSITIES_ALL_PATH)
		if err != nil {
			return
		}
		arrGen(res, &universities)
	}
	//Marks the respons as containting json
	http.Header.Add(w.Header(), "content-type", "application/json")

	//Splits the url from the unisearcher-reqeust into parts
	//If it contains 5 parts, then it is quering uniinfo, and the appropriate func is called
	//if it contains 6 parts, then neighbouringunis is called
	//if the rawquery (i.e. ?limit=) is not empty, that is also forwarded to neighbouring unis
	parts := strings.Split(string(r.URL.Path), "/")
	s := make([]string, 0)

	if len(parts) == 5 {
		writeUniversities(w, findUniversityName(parts[4]))
	}else if len(parts) == 6 {
		s = append(s, parts[4], parts[5])
		if r.URL.RawQuery != "" {
			s = append(s, r.URL.RawQuery)
		}
		HandleCountry(w, s)
	}else {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		log.Println("Malformed URL in request.")
		return
	}
}

//The function that writes the final response for the unisearcher-request
//Used by both findUniversityName and handleCountry
func writeUniversities(w http.ResponseWriter, u []University) {
	encoder := json.NewEncoder(w)
	encoder.Encode(u)
}

//the func that handles neighbourunis-requests
//Fills the unis-struct with all unis that will be sent to user and calls writeUniversites
func HandleCountry(w http.ResponseWriter, s []string) {
	limit := 0
	var unis []University
	var country Country
	var neighbours []Country
	
	//Checks if the country that is being queried exists in the restcountries API
	//Only works if the name being sent is the common name in RESTcountries
	//Russia will work, russian federation will not work
	//LATER: fix this
	for _, val := range countries {
		if val.Name["common"] == strings.Title(s[0]) {
			country = val
			break
		}
	}

	//checks if the final part of the string-slice is the limit param,
	//and if it contains something other 0, that string will be converted to
	//an int that determines how many universities each neighbouring country will respond with
	if len(s) == 3 && s[2] != "" {
		tmp := strings.SplitAfter(s[2], "=")
		_, err := fmt.Sscan(tmp[1], &limit)
		if err != nil {
			limit = 0
		}
	}

	//Tries to find the neighbouring countries for the counry being queried
	//see findNeighbours
	neighbours = append(neighbours, findNeighbours(country.Neighbours)...)

	//For however many neighbours the queried country has,
	//this loop will find either all or a 
	//limit-dependent number of universities that match name
	for _, x := range neighbours {
		if limit == 0 {
			unis = append(unis, findUniversityNameNeighbor(s[1], x)...)
		} else if limit > 0 {
			for i := 0; i < limit; i++ {
				tmp := replySingUniversity(s[1], x)
				//if the country field of the returned university from replySing...is empty
				//then no matching university was found, and the dummy university is not appended
				//as of now, only works if limit is 1, otherwise duplicates will be appended
				if tmp.Country != "" {
					unis = append(unis, tmp)
				}
			}
		}
	}
	
	//adds languages and map-links for the universities
	addMisc(unis)
	writeUniversities(w, unis)
}

//finds all the countries whose cca3-codes are a match with those 
//codes contained in a given country's borders slice
func findNeighbours(codes []string) []Country {
	var c []Country
	for _, x := range countries {
		for _, y := range codes {
			if y == x.Isocode {
				c = append(c, x)
			}
		}
	}
	return c
}

/* ------------------------------------------- */
//The following functions all loop through the universities slice and attempt to find universities
//that match with name/name and country, as is appropriate for the request in question


//replies with a single university that has the appropriate country and name
//if no such university can be found in the universites slice
//an empty object will be returned
func replySingUniversity(name string, country Country) University {
	for _, val := range universities {
		if val.Country == country.Name["common"] && 
		strings.Contains(strings.ToUpper(val.Name), strings.ToUpper(name)) {
			return val
		}
	}
	return University{}
}

//find all universites in the slice that match name either partially or completely
func findUniversityName(name string) []University {
	var tmp []University
	for _, val := range universities {
		if strings.Contains(strings.ToUpper(val.Name), strings.ToUpper(name)) {
			tmp = append(tmp, val)
		}
	}

	//adds languages and map to university
	addMisc(tmp)

	return tmp
}

//finds all universities who are match for both country and name in the universities slice
func findUniversityNameNeighbor(name string, country Country) []University {
	var tmp []University
	for _, val := range universities {
		if val.Country == country.Name["common"] &&
	strings.Contains(strings.ToUpper(val.Name), strings.ToUpper(name)) {
			tmp = append(tmp, val)
		}
	}

	return tmp
}

