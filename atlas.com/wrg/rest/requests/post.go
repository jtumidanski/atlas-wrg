package requests

import (
	"bytes"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PostRequest[A any] func(l logrus.FieldLogger, span opentracing.Span) (DataContainer[A], ErrorListDataContainer, error)

func post(l logrus.FieldLogger, span opentracing.Span) func(url string, input interface{}, resp interface{}, errResp *ErrorListDataContainer) error {
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

func MakePostRequest[A any](url string, i interface{}, configurators ...Configurator) PostRequest[A] {
	return func(l logrus.FieldLogger, span opentracing.Span) (DataContainer[A], ErrorListDataContainer, error) {
		c := &configuration{}
		for _, configurator := range configurators {
			configurator(c)
		}

		r := dataContainer[A]{includedMappers: c.mappers}
		errResp := ErrorListDataContainer{}

		err := post(l, span)(url, i, r, &errResp)
		return r, errResp, err
	}
}
