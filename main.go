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

func buildActorLists(responses []ResponseBody) [][]string {
	actor_list := make([][]string, len(responses))
	for i := 0; i < len(responses); i++ {
		actor_list[i] = strings.Split(responses[i].Actors, ", ")
	}
	return actor_list
}

func compareRequest(urls []url.URL, movie_names []string, reses []ResponseBody) {
	for i := 0; i < len(movie_names); i++ {
		request(urls[i].String(), movie_names[i], &reses[i])
	}
}

func printCommonActors(c_actrs []string) {
	if c_actrs == nil {
		fmt.Println("The films have no actors in common.")
	} else {
		fmt.Println("The films have the following actor(s) in common:")
		for _, value := range c_actrs {
			fmt.Println(value)
		}
	}

}

func findCommon(actor_lists [][]string) []string {
	var c_actors []string
	for i := 0; i < len(actor_lists); i++ {
		compare_list := actor_lists[i]
		stash_map := make(map[string]bool)
		stash_list := make([]string, 0)
		if i+1 < len(actor_lists) {
			compare_list = append(compare_list, actor_lists[i+1]...)
			for _, value := range compare_list {
				if stash_map[value] {
					stash_list = append(stash_list, value)
				}
				if !stash_map[value] {
					stash_map[value] = true
				}
			}
			c_actors = stash_list
			if i+2 < len(actor_lists) {
				new_list_set := [][]string{stash_list, actor_lists[i+2]}
				findCommon(new_list_set)
			}
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
		queries := make([]string, len(films))
		urls := make([]url.URL, len(films))
		r_bodies := make([]ResponseBody, len(films))
		for i := 0; i < len(films); i++ {
			queries[i] = formatQuery(films[i])
			urls[i] = buildUrl(queries[i], "")
			r_bodies[i] = ResponseBody{}
		}
		compareRequest(urls, films, r_bodies)
		common_actors := findCommon(buildActorLists(r_bodies))
		printCommonActors(common_actors)
	default:
		fmt.Println("None is a movie about my fist entering your face.")
	}
}
