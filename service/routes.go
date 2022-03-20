package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-dimension-api/handler"
	"github.com/ONSdigital/dp-cantabular-dimension-api/middleware"

	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/dp-authorisation/auth"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
)

func (svc *Service) buildRoutes(ctx context.Context) {
	svc.router = chi.NewRouter()

	svc.router.Handle("/health", http.HandlerFunc(svc.HealthCheck.Handler))

	if svc.Config.EnablePrivateEndpoints {
		svc.privateEndpoints(ctx)
	} else {
		svc.publicEndpoints(ctx)
	}
}

func (svc *Service) publicEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling public endpoints")

	// Routes
	areaTypes := handler.NewAreaTypes(svc.responder, svc.cantabularClient)
	areas := handler.NewAreas(svc.responder, svc.cantabularClient)
	dimensionOptions := handler.NewDimensionOptions(svc.responder, svc.cantabularClient)

	svc.router.Get("/area-types", areaTypes.Get)
	svc.router.Get("/areas", areas.Get)
	svc.router.Get("/dimension/{dimension}/options", dimensionOptions.Get)
}

func (svc *Service) privateEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling private endpoints")

	// use sub-router to prevent /health endpoint from requiring auth
	r := chi.NewRouter()

	// Middleware
	svc.identityClient = identity.New(svc.Config.ZebedeeURL)
	checkIdentity := dphandlers.IdentityWithHTTPClient(svc.identityClient)
	permissions := middleware.NewPermissions(svc.Config.ZebedeeURL, svc.Config.EnablePermissionsAuth)

	r.Use(permissions.Require(auth.Permissions{Read: true}))
	r.Use(checkIdentity)
	r.Use(middleware.LogIdentity())

	// Routes
	areaTypes := handler.NewAreaTypes(svc.responder, svc.cantabularClient)
	areas := handler.NewAreas(svc.responder, svc.cantabularClient)
	dimensionOptions := handler.NewDimensionOptions(svc.responder, svc.cantabularClient)

	r.Get("/area-types", areaTypes.Get)
	r.Get("/areas", areas.Get)
	r.Get("/dimension/{dimension}/options", dimensionOptions.Get)

	svc.router.Mount("/", r)
}
