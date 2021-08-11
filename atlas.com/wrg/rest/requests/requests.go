package requests

import (
	json2 "atlas-wrg/json"
	"atlas-wrg/retry"
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	BaseRequest string = "http://atlas-nginx:80"
)

type configuration struct {
	retries int
}

type Configurator func(c *configuration)

func SetRetries(amount int) Configurator {
	return func(c *configuration) {
		c.retries = amount
	}
}

func Get(l logrus.FieldLogger) func(url string, resp interface{}, configurators ...Configurator) error {
	return func(url string, resp interface{}, configurators ...Configurator) error {
		c := &configuration{retries: 1}
		for _, configurator := range configurators {
			configurator(c)
		}

		var r *http.Response
		get := func(attempt int) (bool, error) {
			var err error
			r, err = http.Get(url)
			if err != nil {
				l.Warnf("Failed calling GET on %s, will retry.", url)
				return true, err
			}
			return false, nil
		}
		err := retry.Try(get, c.retries)
		if err != nil {
			l.WithError(err).Errorf("Unable to successfully call GET on %s.", url)
			return err
		}
		err = processResponse(r, resp)
		return err
	}
}

func Post(url string, input interface{}) (*http.Response, error) {
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

func Delete(url string) (*http.Response, error) {
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
