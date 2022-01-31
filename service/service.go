package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-dimension-api/config"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config           *config.Config
	Server           HTTPServer
	router           *chi.Mux
	responder        Responder
	cantabularClient CantabularClient
	HealthCheck      HealthChecker
	identityClient   *identity.Client
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Init(ctx context.Context, buildTime, gitCommit, version string) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	log.Info(ctx, "initialising service with config", log.Data{"config": cfg})

	svc.Config = cfg

	svc.HealthCheck, err = GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return fmt.Errorf("failed to get healthcheck: %w", err)
	}

	svc.responder = GetResponder()
	svc.cantabularClient = GetCantabularClient(cfg)

	svc.buildRoutes(ctx)
	svc.Server = GetHTTPServer(cfg.BindAddr, svc.router)

	if err := svc.registerCheckers(); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	svc.HealthCheck.Start(ctx)

	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failed to start main http server: %w", err)
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.HealthCheck != nil {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
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

func (svc *Service) registerCheckers() error {
	hc := svc.HealthCheck

	// TODO - when Cantabular server is deployed to Production, remove this placeholder and the flag,
	// and always use the real Checker instead: svc.cantabularClient.Checker
	cantabularChecker := svc.cantabularClient.Checker
	cantabularAPIExtChecker := svc.cantabularClient.CheckerAPIExt
	if !svc.Config.CantabularHealthcheckEnabled {
		cantabularChecker = func(ctx context.Context, state *healthcheck.CheckState) error {
			return state.Update(healthcheck.StatusOK, "Cantabular healthcheck placeholder", http.StatusOK)
		}
		cantabularAPIExtChecker = func(ctx context.Context, state *healthcheck.CheckState) error {
			return state.Update(healthcheck.StatusOK, "Cantabular APIExt healthcheck placeholder", http.StatusOK)
		}
	}

	if err := hc.AddCheck("Cantabular server", cantabularChecker); err != nil {
		return fmt.Errorf("error adding check for Cantabular server: %w", err)
	}

	if err := hc.AddCheck("Cantabular API Extension", cantabularAPIExtChecker); err != nil {
		return fmt.Errorf("error adding check for Cantabular api extension: %w", err)
	}

	return nil
}
