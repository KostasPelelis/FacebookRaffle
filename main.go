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
			post := data[i].(map[string]interface{})
			fmt.Println(post)
			joined_users = append(joined_users, post["from"].(string))
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
	if val, ok := (*ops)["winners"]; ok {
		lim, _ := strconv.Atoi(val)
		for i := 0; i < lim; i++ {
			winner_index := rand.Intn(len(joined_users))
			// Here we should check the option liked_users_only
			// That can be done using the facebook api to cross
			// validate the potential winner if he also likes the page
			winners = append(winners, joined_users[winner_index])
			if _, okw := (*ops)["unique_winners"]; okw {
				joined_users = append(joined_users[:winner_index], joined_users[winner_index:]...)
			}
		}
	}
}

var available_options []string = []string{"user_id", "post_id", "winners", "liked_users_only", "unique_winners"}
var api_base string = "https://graph.facebook.com"

func WebRaffleHandler(wr http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.RawQuery)
	var options map[string]string
	options = make(map[string]string)
	q, _ := url.ParseQuery(req.URL.RawQuery)
	for i := range available_options {
		option := available_options[i]
		if val, ok := q[option]; ok {
			switch option {
			case "user_id":
				options["user_id"] = val[0]
			case "post_id":
				options["post_id"] = val[0]
			case "winners":
				options["winners"] = val[0]
			case "liked_users_only":
				options["liked_users_only"] = val[0]
			case "unique_winners":
				options["unique_winners"] = val[0]
			default:
				panic("Unrecognized option")
			}
		}
	}
	Raffle(api_base, options["post_id"], options["user_id"], &options)
}

func main() {
	http.HandleFunc("/raffle", WebRaffleHandler)
	panic(http.ListenAndServe(":8080", nil))
}
