package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ResponseBody struct {
	Actors string
}

var search = flag.String("s", "None", "search for a movie or series")
var base_query url.URL

// = "http://omdbapi.com/?r=json&plot=short&t="

func main() {
	base_query.Scheme = "http"
	base_query.Host = "omdbapi.com"
	the_query := base_query.Query()
	the_query.Set("r", "json")
	the_query.Add("plot", "short")
	flag.Parse()
	query := strings.Join(strings.Split(strings.ToLower(*search), " "), "+")
	the_query.Add("t", query)
	base_query.RawQuery = the_query.Encode()

	resp, err := http.Get(base_query.String())
	if err != nil {
		log.Fatal(err)
		fmt.Println("The following movie was not found: ", *search)
	}
	if resp.StatusCode != 200 {
		log.Fatal(err)
	}
	actor_list := ResponseBody{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&actor_list)
	fmt.Println(actor_list)
	for _, value := range strings.Split(actor_list.Actors, ", ") {
		fmt.Println(value)
	}
}
