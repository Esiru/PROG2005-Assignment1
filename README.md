# Assignment 1, PROG2005, spring semester 2023
# Author: Even Stetrud

An attempt at writing a RESTful API (unisearcher) that fetches information from two other APIs (http://universities.hipolabs.com/ and https://restcountries.com/) and presents information about universities, including name, country, web-domains and languages. 

## Deployment
Deployed on https://unisearcher-k43m.onrender.com/unisearcher/v1/neighbourunis/ . There are three end-points:
neighbourunis/<country>/<name>?limit=<limit>
uniinfo/<name>
/diag/

Neighbourunis provides information about universities in the neighbouring countries of <countyr>, and whose names are a partial or complete names with <name>. ?limit is an optional parameter that limits the number of universities that will be responded, PER  NEIGHBOURING COUNTRY.
uniinfo provides information about all universities whose names are a partial or complete match with <name>
diag is a status endpoint that displays the status code of the two APIs, the current version of unisearcher and the uptime of the current deployment

## to-do
 - fix functionality for those universites whose country name in hipolabs-univeristy API do not match the common-name of the correspoding country in the RESTcountries API. (E.g. Russian universities are in the country Russian Federation in the hipolabs API, while the common name of that country is Russia in RESTcountries)