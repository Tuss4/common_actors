package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ResponseBody struct {
	Actors string
}

var search = flag.String("s", "None", "search for a movie or series")
var base_query = "http://omdbapi.com/?r=json&plot=short&t="

func main() {
	flag.Parse()
	fmt.Println("Search has value ", *search)
	query := strings.Join(strings.Split(strings.ToLower(*search), " "), "+")
	url := base_query + query
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(err)
	}
	actor_list := ResponseBody{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&actor_list)
	fmt.Println(actor_list)
}