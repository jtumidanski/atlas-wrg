package requests

import (
	json2 "atlas-wrg/json"
	"atlas-wrg/retry"
	"bytes"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
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

func Get(l logrus.FieldLogger, span opentracing.Span) func(url string, resp interface{}, configurators ...Configurator) error {
	return func(url string, resp interface{}, configurators ...Configurator) error {
		c := &configuration{retries: 1}
		for _, configurator := range configurators {
			configurator(c)
		}

		var r *http.Response
		get := func(attempt int) (bool, error) {
			var err error

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				l.WithError(err).Errorf("Error creating request.")
				return true, err
			}
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			err = opentracing.GlobalTracer().Inject(
				span.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(req.Header))
			if err != nil {
				l.WithError(err).Errorf("Unable to decorate request headers with OpenTracing information.")
			}
			r, err = http.DefaultClient.Do(req)
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
		err = ProcessResponse(r, resp)
		return err
	}
}

func Post(l logrus.FieldLogger, span opentracing.Span) func(url string, input interface{}) (*http.Response, error) {
	return func(url string, input interface{}) (*http.Response, error) {
		jsonReq, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonReq))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		err = opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			l.WithError(err).Errorf("Unable to decorate request headers with OpenTracing information.")
		}
		return http.DefaultClient.Do(req)
	}
}

func ProcessResponse(r *http.Response, rb interface{}) error {
	err := json2.FromJSON(rb, r.Body)
	if err != nil {
		return err
	}

	return nil
}

func ProcessErrorResponse(r *http.Response, eb interface{}) error {
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
