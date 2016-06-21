package apiwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type API struct {
	base_url string
	client   http.Client
}

func NewAPI(base_url string) *API {
	hc := http.Client{}
	return &API{base_url, hc}
}

func (api API) Send(req *http.Request) *http.Response {
	resp, _ := api.client.Do(req)
	return resp
}

func (api API) MakeRequest(url string, method string, params map[string]string) *http.Request {
	req, err := http.NewRequest(method, api.base_url+url, nil)
	if err != nil {
		fmt.Println("Something went wrong while building HTTP Request %s", err)
		log.Fatal(err)
		return nil
	}
	q := req.URL.Query()
	if params != nil {
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	fmt.Printf("[%s] %s%s?%s\n", method, api.base_url, url, req.URL.RawQuery)
	return req
}

func (api API) Get(url string, params map[string]string) map[string]interface{} {
	get_req := api.MakeRequest(url, "GET", params)
	res := api.Send(get_req)
	if res == nil {
		return nil
	}
	resBuf := bytes.Buffer{}
	defer res.Body.Close()
	_, err := io.Copy(&resBuf, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data map[string]interface{}
	if errJ := json.Unmarshal(resBuf.Bytes(), &data); errJ != nil {
		log.Fatal(err)
	}
	return data
}
