package RESTUniversity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

//the func that sends http-requests for all the queries, more or less a duplicate of the 
//client code from 05-REST-client
func Client(url string) (*http.Response, error) {

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
		return nil, err
	}

	r.Header.Add("content-type", "application/json")

	client := &http.Client{}
	defer client.CloseIdleConnections()

	res, err := client.Do(r)
	if err != nil {
		fmt.Errorf("Error in response: ", err.Error())
		return nil, err
	}

	//log.Println("URL:", url, "Content length:", res.ContentLength)
	if res.ContentLength == 2 {
		return nil, errors.New("no resource found")
	}

	fmt.Println("Status:", res.Status)
	fmt.Println("Status code:", res.StatusCode)

	fmt.Println("Content type:", res.Header.Get("Content-type"))
	fmt.Println("Protocol:", res.Proto)
	return res, nil
}

//a generic function that unmarshals the contents of a response body
//into either a university or country slice
func arrGen[v *[]University | *[]Country](res *http.Response, arr v) (error) {

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	res.Body.Close()

	if err := json.Unmarshal([]byte(body), &arr); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//finds the languages and maps for a university, as decided by which country it belongs to
//does this by looping through the countries slice for EACH university in the passed slice/array
func addMisc[v []University | *[1]University](arr v) (error) {
	for i := 0; i < len(arr); i++ {
		 e := &arr[i]
		for j := range countries {
			if e.Country == countries[j].Name["common"] {	
				e.Languages = countries[j].Languages
				e.Map = countries[j].Maps
				break
			}
		}
	}
	return nil
}