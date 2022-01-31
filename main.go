package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-cantabular-dimension-api/config"
	"github.com/ONSdigital/dp-cantabular-dimension-api/service"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-cantabular-dimension-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	svcErrors := make(chan error, 1)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("unable to retrieve configuration, error: %w", err)
	}
	log.Info(ctx, "config on startup", log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Run the service
	svc := service.New()
	if err := svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}
	svc.Start(ctx, svcErrors)

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		err = fmt.Errorf("service error received: %w", err)
		if errClose := svc.Close(ctx); errClose != nil {
			log.Error(ctx, "service close error during error handling", errClose)
		}
		return err
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}

	if err := svc.Close(ctx); err != nil {
		return fmt.Errorf("failed to close service: %w", err)
	}

	return nil
}
