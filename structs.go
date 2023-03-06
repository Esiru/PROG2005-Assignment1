package RESTUniversity

type Diag struct {
	UniStatus 		int     `json:"universitiesapi"`
	CountryStatus   int     `json:"countriesapi"`
	Version 		string  `json:"version"`
	Uptime 			float64 `json:"uptime"`
}
	
type Country struct {
	Neighbours []string `json:"borders"`
	Isocode    string `json:"cca3"`
	Name	   map[string]interface{}
	Languages  map[string]string
	Maps	   map[string]string 
}

type University struct {
	Name 	  string `json:"name"`
	Country   string `json:"country"`
	Isocode   string `json:"alpha_two_code"`
	Web_pages []string `json:"web_pages"`
	Languages map[string]string `json:"languages"`
	Map		  map[string]string `json:"maps"`
}
