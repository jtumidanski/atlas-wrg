package retry

import (
	"errors"
	"time"
)

type RetryFunc func(attempt int) (retry bool, err error)

type RetryResponseFunc func(attempt int) (bool, interface{}, error)

func Retry(fn RetryFunc, retries int) error {
	attempt := 1
	for {
		cont, err := fn(attempt)
		if !cont || err == nil {
			return nil
		}
		attempt++
		if attempt > retries {
			return errors.New("max retry reached")
		}
		time.Sleep(1 * time.Second)
	}
}

func RetryResponse(fn RetryResponseFunc, retries int) (interface{}, error) {
	attempt := 1
	for {
		cont, m, err := fn(attempt)
		if !cont || err == nil {
			return m, nil
		}
		attempt++
		if attempt > retries {
			return nil, errors.New("max retry reached")
		}
		time.Sleep(1 * time.Second)
	}
}