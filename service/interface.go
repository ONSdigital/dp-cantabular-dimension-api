package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-dimension-api/config"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mock/initialiser.go -pkg mock . Initialiser
//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/healthCheck.go -pkg mock . HealthChecker

// Initialiser defines the methods to initialise external services
type Initialiser interface {
	DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer
	DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error)
}

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddCheck(name string, checker healthcheck.Checker) (err error)
}

// Responder handles responding to http requests
type Responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
	StatusCode(http.ResponseWriter, int)
	Bytes(context.Context, http.ResponseWriter, int, []byte)
}

type CantabularClient interface {
	GetCodebook(context.Context, cantabular.GetCodebookRequest) (*cantabular.GetCodebookResponse, error)
	GetGeographyDimensions(context.Context, cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)
	GetAreas(context.Context, cantabular.QueryData) (*cantabular.GetAreasResponse, error)
	Checker(context.Context, *healthcheck.CheckState) error
	CheckerAPIExt(ctx context.Context, state *healthcheck.CheckState) error
	StatusCode(error) int
}
