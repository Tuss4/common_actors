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
var single_list, list_one, list_two ResponseBody
var base_query url.URL

func build_url(query string) url.URL {
	base_query.Scheme = "http"
	base_query.Host = "omdbapi.com"
	the_query := base_query.Query()
	the_query.Set("r", "json")
	the_query.Add("plot", "short")
	the_query.Add("t", query)
	base_query.RawQuery = the_query.Encode()
	return base_query
}

func format_query(q string) string {
	return strings.Join(strings.Split(strings.ToLower(q), " "), "+")
}

func movie_not_found(movie_name string) {
	fmt.Println("The following movie was not found:", movie_name)
}

func handle_error_n_status_code(err error, respsonse *http.Response, movie_name string) {
	if err != nil {
		log.Fatal(err)
		movie_not_found(movie_name)
	}
	if respsonse.StatusCode != 200 {
		log.Fatal(err)
	}
}

func request(url string, movie_name string, res *ResponseBody) {
	resp, err := http.Get(url)
	handle_error_n_status_code(err, resp, movie_name)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(res)
}

func compare_request(url_1 string, url_2 string, movie_names []string, res_1 *ResponseBody, res_2 *ResponseBody) {
	resp, err := http.Get(url_1)
	handle_error_n_status_code(err, resp, movie_names[0])
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(res_1)

	resp_2, err_2 := http.Get(url_2)
	handle_error_n_status_code(err_2, resp_2, movie_names[1])
	defer resp_2.Body.Close()
	err_2 = json.NewDecoder(resp_2.Body).Decode(res_2)

}
func find_common(l1, l2 []string) map[string]bool {
	actor_map := make(map[string]bool)
	for _, value := range l1 {
		actor_map[value] = true
	}
	return actor_map

}

func main() {
	flag.Parse()
	query := format_query(*search)
	if *search != "" {
		url := build_url(query)
		request(url.String(), *search, &single_list)
		for _, value := range strings.Split(single_list.Actors, ", ") {
			fmt.Println(value)
		}
	} else if *common != "" {
		films := strings.Split(*common, ", ")
		fmt.Println("You're trying to find common actors between:", films[0], "and", films[1])
		query_1 := format_query(films[0])
		url_1 := build_url(query_1)
		query_2 := format_query(films[1])
		url_2 := build_url(query_2)
		compare_request(
			url_1.String(), url_2.String(), films, &list_one, &list_two)
		fmt.Println(list_two.Actors)
	} else {
		fmt.Println("None is a movie about my fist entering your face.")
	}
}
