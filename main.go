package main

import (
	"./apiwrapper"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func Raffle(base string, post_id string, user_id string, ops *map[string]string) {

	fb := apiwrapper.NewAPI(base)

	var params map[string]string
	params = make(map[string]string)

	params["access_token"] = os.Getenv("FB_ACCESS_TOKEN")

	current_url := "/v2.6/" + user_id + "_" + post_id + "/sharedposts"

	var joined_users []string = []string{}
	for current_url != "" {
		res := fb.Get(current_url, params)
		data := res["data"].([]interface{})
		for i := range data {
			fmt.Println(i)
			//post := data[i].(map[string]interface{})
			joined_users = append(joined_users, "asc")
		}
		paging := res["paging"].(map[string]interface{})
		u, err := url.Parse(paging["next"].(string))
		if err != nil {
			panic(err)
		}
		qp, _ := url.ParseQuery(u.RawQuery)
		for k, v := range qp {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
		current_url = u.Path
	}
	var winners []string = []string{}
	if val, ok := (*ops)["WINNERS"]; ok {
		lim, _ := strconv.Atoi(val)
		for i := 0; i < lim; i++ {
			winner_index := rand.Intn(len(joined_users))
			winners = append(winners, joined_users[winner_index])
			if _, okw := (*ops)["UNIQUE_WINNERS"]; okw {
				joined_users = append(joined_users[:winner_index], joined_users[winner_index:]...)
			}
		}
	}
}

var available_options []string = []string{"USER_ID", "POST_ID", "WINNERS", "LIKED_USERS_ONLY", "UNIQUE_WINNERS"}

func WebRaffleHandler(wr http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.RawQuery)
	api_base := "https://graph.facebook.com"
	var options map[string]string
	options = make(map[string]string)
	q, _ := url.ParseQuery(req.URL.RawQuery)
	for i := range available_options {
		option := available_options[i]
		if val, ok := q[option]; ok {
			switch option {
			case "USER_ID":
				options["USER_ID"] = val[0]
			case "POST_ID":
				options["POST_ID"] = val[0]
			case "WINNERS":
				options["WINNERS"] = val[0]
			case "LIKED_USERS_ONLY":
				options["LIKED_USERS_ONLY"] = val[0]
			case "UNIQUE_WINNERS":
				options["UNIQUE_WINNERS"] = val[0]
			default:
				panic("Unrecognized option")
			}
		}
	}
	fmt.Println(options)
	Raffle(api_base, options["POST_ID"], options["USER_ID"], &options)
}

func main() {
	fmt.Println("Starting")
	http.HandleFunc("/raffle", WebRaffleHandler)
	panic(http.ListenAndServe(":8080", nil))
}
