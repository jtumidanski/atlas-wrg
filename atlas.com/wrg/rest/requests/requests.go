package requests

import (
	json2 "atlas-wrg/json"
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	BaseRequest string = "http://atlas-nginx:80"
)

func get(url string, resp interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}

	err = processResponse(r, resp)
	return err
}

func post(url string, input interface{}) (*http.Response, error) {
	jsonReq, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	r, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func delete(url string) (*http.Response, error) {
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	return client.Do(r)
}

func processResponse(r *http.Response, rb interface{}) error {
	err := json2.FromJSON(rb, r.Body)
	if err != nil {
		return err
	}

	return nil
}

func processErrorResponse(r *http.Response, eb interface{}) error {
	if r.ContentLength > 0 {
		err := json2.FromJSON(eb, r.Body)
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}
