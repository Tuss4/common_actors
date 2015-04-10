package main

import (
	"flag"
	"fmt"
	// "net/http"
	"strings"
	//"log"
)

var search = flag.String("s", "None", "search for a movie or series")
var url = "http://omdbapi.com/?r=json&plot=short&t="

func main() {
	flag.Parse()
	fmt.Println("Search has value ", *search)
	query := strings.Join(strings.Split(strings.ToLower(*search), " "), "+")
	fmt.Println(query)
}
