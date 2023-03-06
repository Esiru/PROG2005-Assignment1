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

func HandlerUniversity() func(http.ResponseWriter, *http.Request) {
	//log.Println("You are here, start of handler-func")
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			//log.Println("MEthod is get")
			handleUniversityGet(w, r/*, db*/)
		default:
			http.Error(w, "Method not allowed, only GET requests are supported", http.StatusNotImplemented)
		}
	}
}

func handleUniversityGet(w http.ResponseWriter, r *http.Request) {
	if len(countries) == 0 {
		res, err := Client(COUNTRIES_ALL_PATH)
		if err != nil {
			log.Println("something went wrong with countries", err)
			return
		}
		arrGen(res, &countries)
	}
	if len(universities) == 0 {
		res, err := Client(UNIVERSITIES_ALL_PATH)
		if err != nil {
			log.Println("something went wrong with unis", err)
			return
		}
		arrGen(res, &universities)
	}
	http.Header.Add(w.Header(), "content-type", "application/json")

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

func writeUniversities(w http.ResponseWriter, u []University) {
	encoder := json.NewEncoder(w)
	encoder.Encode(u)
}

func HandleCountry(w http.ResponseWriter, s []string) {
	limit := 0
	var unis []University
	var country Country
	var neighbours []Country
	
	for _, val := range countries {
		if val.Name["common"] == strings.Title(s[0]) {
			country = val
			break
		}
	}

	if len(s) == 3 && s[2] != "" {
		tmp := strings.SplitAfter(s[2], "=")
		_, err := fmt.Sscan(tmp[1], &limit)
		if err != nil {
			limit = 0
		}
	}

	neighbours = append(neighbours, findNeighbours(country.Neighbours)...)

	log.Println(country)
	log.Println(neighbours)
	for _, x := range neighbours {
		if limit == 0 {
			unis = append(unis, findUniversityNameNeighbor(s[1], x)...)
		} else if limit > 0 {
			for i := 0; i < limit; i++ {
				unis = append(unis, replySingUniversity(s[1], x))	
			}
		}
	}
	
	addMisc(unis)
	writeUniversities(w, unis)
}

func replySingUniversity(name string, country Country) University {
	var tmp University
	for _, val := range universities {
		if val.Name ==  country.Name["common"] && 
		strings.Contains(strings.ToUpper(val.Name), strings.ToUpper(name)) {
			return tmp
		}
	}
	return University{}
}

func findUniversityName(name string) []University {
	var tmp []University
	for _, val := range universities {
		if strings.Contains(strings.ToUpper(val.Name), strings.ToUpper(name)) {
			log.Println(val)
			tmp = append(tmp, val)
		}
	}

	addMisc(tmp)

	return tmp
}

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