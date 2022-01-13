package service

import (
	"context"
	"fmt"
	"errors"

	"github.com/ONSdigital/dp-cantabular-dimension-api/config"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	config      *config.Config
	server      HTTPServer
	router      *chi.Mux
	responder   Responder
	healthCheck HealthChecker
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Init(ctx context.Context, buildTime, gitCommit, version string) error {
	cfg, err := config.Get()
	if err != nil{
		return fmt.Errorf("failed to get config: %w", err)
	}

	log.Info(ctx, "initialising service with config", log.Data{"config": cfg})

	svc.config = cfg

	svc.healthCheck, err = GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return fmt.Errorf("failed to get healthcheck: %w", err)
	}

	svc.responder = GetResponder()

	svc.buildRoutes()
	svc.server = GetHTTPServer(cfg.BindAddr, svc.router)

	if err := registerCheckers(ctx, svc.healthCheck); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	svc.healthCheck.Start(ctx)

	go func() {
		if err := svc.server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failed to start main http server: %w", err)
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.healthCheck != nil {
			svc.healthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}

	// other error
	if hasShutdownError {
		return errors.New("failed to shutdown gracefully")
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context, hc HealthChecker) error {
	return nil
}
