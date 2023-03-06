package RESTUniversity

//Default port
const DEFAULT_PORT = "8080"

//remaining consts contain complete or partial URLs for unisearcher and the two attendant APIs
//some may be obsolete
const UNISEARCHER_PATH = "/unisearcher/v1/"

const UNIVERSITIES_PATH = "http://universities.hipolabs.com/search?name="

const COUNTRY_SUFFIX = "&country="

const COUNTRIES_NAME_PATH = "https://restcountries.com/v3.1/name/"

const FULLTEXT_SUFFIX = "fullText=true"

const COUNTRIES_CODE_PATH = "https://restcountries.com/v3.1/alpha?codes="

const COUNTRIES_ALL_PATH = "https://restcountries.com/v3.1/all?fields=name,cca3,borders,languages,maps"

const UNIVERSITIES_ALL_PATH = "http://universities.hipolabs.com/search?name="

const LIMIT_PARAM = "?limit="