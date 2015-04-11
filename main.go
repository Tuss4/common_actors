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

var search = flag.String("s", "", "search for a movie or series")
var common = flag.String("c", "", "find common actors")

// var api_key = flag.String("k", "", "api key")
var base_query url.URL

func build_url(q string) url.URL {
	base_query.Scheme = "http"
	base_query.Host = "omdbapi.com"
	the_query := base_query.Query()
	the_query.Set("r", "json")
	the_query.Add("plot", "short")
	the_query.Add("t", q)
	base_query.RawQuery = the_query.Encode()
	return base_query
}

func request(u string, q string, res *ResponseBody) {
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
		fmt.Println("The following movie was not found: ", q)
	}
	if resp.StatusCode != 200 {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(res)
}

func main() {
	flag.Parse()
	query := strings.Join(strings.Split(strings.ToLower(*search), " "), "+")
	fmt.Println(query)
	if *search != "" {
		url := build_url(query)
		actor_list := ResponseBody{}
		request(url.String(), *search, &actor_list)
		fmt.Println(actor_list)
		for _, value := range strings.Split(actor_list.Actors, ", ") {
			fmt.Println(value)
		}
	} else if *common != "" {
		films := strings.Split(*common, ", ")
		fmt.Println(films)
		fmt.Println("You're trying to find common actors between:", films[0], "and", films[1])
	} else {
		fmt.Println("None is a movie about my fist entering your face.")
	}
}
