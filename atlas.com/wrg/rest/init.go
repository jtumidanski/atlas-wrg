package rest

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

func CreateRestService(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg)
}
