package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	gokeymimepi "github.com/kai5263499/go-key-mime-pi"
	"github.com/kai5263499/go-key-mime-pi/internal/api"
	"github.com/kai5263499/go-key-mime-pi/internal/domain"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := &domain.Config{}
	if err := env.Parse(cfg); err != nil {
		log.WithError(err).Fatal("parse config")
	}

	if level, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("parse log level")
	} else {
		log.SetLevel(level)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetReportCaller(true)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	h, err := gokeymimepi.NewHid()
	if err != nil {
		log.WithError(err).Fatal("new hid")
	}

	api, err := api.New(
		ctx,
		stop,
		cfg,
		h,
	)
	if err != nil {
		log.WithError(err).Fatal("new api")
	}

	api.Start()

	<-ctx.Done()

	log.Info("api exiting")
}
