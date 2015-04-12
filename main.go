package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	// "sync"
)

type ResponseBody struct {
	Actors string
}

var search = flag.String("s", "", "search for a movie or series")
var common = flag.String("c", "", "find common actors")
var year = flag.String("y", "", "specify a specific year")
var single_list, list_one, list_two ResponseBody

func buildUrl(query, year string) url.URL {
	var base_query url.URL
	base_query.Scheme = "http"
	base_query.Host = "omdbapi.com"
	the_query := base_query.Query()
	the_query.Set("r", "json")
	the_query.Add("plot", "short")
	the_query.Add("t", query)
	if year != "" {
		the_query.Add("y", year)
	}
	base_query.RawQuery = the_query.Encode()
	return base_query
}

func formatQuery(q string) string {
	return strings.Join(strings.Split(strings.ToLower(q), " "), "+")
}

func movieNotFound(movie_name string) {
	fmt.Println("The following movie was not found:", movie_name)
}

func handleErrorNStatusCode(err error, respsonse *http.Response, movie_name string) {
	if err != nil {
		log.Fatal(err)
		movieNotFound(movie_name)
	}
	if respsonse.StatusCode != 200 {
		log.Fatal(err)
	}
}

func request(url string, movie_name string, res *ResponseBody) {
	resp, err := http.Get(url)
	handleErrorNStatusCode(err, resp, movie_name)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(res)
}

func compareRequest(url_1 string, url_2 string, movie_names []string, res_1 *ResponseBody, res_2 *ResponseBody) {
	request(url_1, movie_names[0], res_1)
	request(url_2, movie_names[1], res_2)
	// var wg sync.WaitGroup

	// wg.Add(2)

	// go func() {
	//        request(url_1, movie_names[0], res_1)
	//        wg.Done()
	//    }

	// go func () {
	//        request(url_2, movie_names[1], res_2)
	//        wg.Done()
	//    }

	//    wg.Wait()
}

func findCommon(l1, l2 string) []string {
	sl1 := strings.Split(l1, ", ")
	sl2 := strings.Split(l2, ", ")
	sl1 = append(sl1, sl2...)
	fmt.Println(sl1)
	c_actors := make([]string, 0)
	actor_map := make(map[string]bool)
	for _, value := range sl1 {
		if actor_map[value] {
			c_actors = append(c_actors, value)
		}
		if !actor_map[value] {
			actor_map[value] = true
		}
	}
	return c_actors

}

func main() {
	flag.Parse()
	query := formatQuery(*search)
	switch {
	case *search != "":
		url := buildUrl(query, *year)
		request(url.String(), *search, &single_list)
		for _, value := range strings.Split(single_list.Actors, ", ") {
			fmt.Println(value)
		}
	case *common != "":
		films := strings.Split(*common, ", ")
		query_1 := formatQuery(films[0])
		url_1 := buildUrl(query_1, "")
		query_2 := formatQuery(films[1])
		url_2 := buildUrl(query_2, "")
		compareRequest(
			url_1.String(), url_2.String(), films, &list_one, &list_two)
		common_actors := findCommon(list_one.Actors, list_two.Actors)
		fmt.Println(films[0], "and", films[1], "have the following actor(s) in common:")
		for _, value := range common_actors {
			fmt.Println(value)
		}
	default:
		fmt.Println("None is a movie about my fist entering your face.")
	}
}
