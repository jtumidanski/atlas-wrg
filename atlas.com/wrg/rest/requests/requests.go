package requests

import (
	json2 "atlas-wrg/json"
	"atlas-wrg/rest/response"
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
		err = processResponse(r, resp)

		l.WithFields(logrus.Fields{"method": http.MethodGet, "status": r.Status, "path": url, "response": resp}).Debugf("Printing request.")

		return err
	}
}

type ErrorListDataContainer struct {
	Errors []ErrorData `json:"errors"`
}

type ErrorData struct {
	Status int               `json:"status"`
	Code   string            `json:"code"`
	Title  string            `json:"title"`
	Detail string            `json:"detail"`
	Meta   map[string]string `json:"meta"`
}

func Post(l logrus.FieldLogger, span opentracing.Span) func(url string, input interface{}, resp interface{}, errResp *ErrorListDataContainer) error {
	return func(url string, input interface{}, resp interface{}, errResp *ErrorListDataContainer) error {
		jsonReq, err := json.Marshal(input)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonReq))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		err = opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			l.WithError(err).Errorf("Unable to decorate request headers with OpenTracing information.")
		}
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if r.StatusCode != http.StatusNoContent && r.StatusCode != http.StatusCreated && r.StatusCode != http.StatusAccepted {
			err = processErrorResponse(r, errResp)
			if err != nil {
				return err
			}

			l.WithFields(logrus.Fields{"method": http.MethodPost, "status": r.Status, "path": url, "input": input, "response": errResp}).Debugf("Printing request.")

			return nil
		}

		if r.ContentLength > 0 {
			err = processResponse(r, resp)
			if err != nil {
				return err
			}
			l.WithFields(logrus.Fields{"method": http.MethodPost, "status": r.Status, "path": url, "input": input, "response": resp}).Debugf("Printing request.")
		} else {
			l.WithFields(logrus.Fields{"method": http.MethodPost, "status": r.Status, "path": url, "input": input, "response": ""}).Debugf("Printing request.")
		}

		return nil
	}
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

type DataContainer[A any] interface {
	Data() DataBody[A]
	DataList() []DataBody[A]
	Length() int
}

type dataContainer[A any] struct {
	data     response.DataSegment
	included response.DataSegment
}

type DataBody[A any] struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes A      `json:"attributes"`
}

func (c *dataContainer[A]) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(c.data) == 1 {
		t.Data = c.data[0]
	} else {
		t.Data = c.data
	}
	return json.Marshal(t)
}

func (c *dataContainer[A]) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyDataBody[A]))
	if err != nil {
		return err
	}

	c.data = d
	c.included = i
	return nil
}

func (c dataContainer[A]) Data() DataBody[A] {
	if len(c.data) >= 1 {
		return *c.data[0].(*DataBody[A])
	}
	return DataBody[A]{}
}

func (c dataContainer[A]) DataList() []DataBody[A] {
	var r = make([]DataBody[A], 0)
	for _, x := range c.data {
		r = append(r, *x.(*DataBody[A]))
	}
	return r
}

func (c dataContainer[A]) Length() int {
	return len(c.data)
}

func EmptyDataBody[A any]() interface{} {
	return &DataBody[A]{}
}

type Request[A any] func(l logrus.FieldLogger, span opentracing.Span) (DataContainer[A], error)

type PostRequest[A any] func(l logrus.FieldLogger, span opentracing.Span) (DataContainer[A], ErrorListDataContainer, error)

func MakeGetRequest[A any](url string, configurators ...Configurator) Request[A] {
	return func(l logrus.FieldLogger, span opentracing.Span) (DataContainer[A], error) {
		r := dataContainer[A]{}
		err := Get(l, span)(url, &r, configurators...)
		return r, err
	}
}

func MakePostRequest[A any](url string, i interface{}) PostRequest[A] {
	return func(l logrus.FieldLogger, span opentracing.Span) (DataContainer[A], ErrorListDataContainer, error) {
		r := dataContainer[A]{}
		errResp := ErrorListDataContainer{}

		err := Post(l, span)(url, i, r, &errResp)
		return r, errResp, err
	}
}
