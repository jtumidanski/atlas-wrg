package main

import (
	"atlas-wrg/channel"
	"atlas-wrg/kafka/consumers"
	"atlas-wrg/logger"
	"atlas-wrg/rest"
	"atlas-wrg/tracing"
	"atlas-wrg/world"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-wrg"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	consumers.CreateEventConsumers(l, ctx, wg)

	rest.CreateService(l, ctx, wg, "/ms/wrg", channel.InitResource, world.InitResource)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
