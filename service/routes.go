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

func (svc *Service) buildRoutes(ctx context.Context){
	if svc.config.EnablePrivateEndpoints{
		svc.router = svc.privateEndpoints(ctx)
	} else{
		svc.router = svc.publicEndpoints(ctx)
	}

	svc.router.Handle("/health", http.HandlerFunc(svc.healthCheck.Handler))
}

func (svc *Service) publicEndpoints(ctx context.Context) *chi.Mux {
	log.Info(ctx, "enabling public endpoints")
	
	r := chi.NewRouter()

	hello := handler.NewHello(svc.responder, svc.cantabularClient)
	r.Get("/hello", hello.Get)
	r.Post("/hello", hello.Create)

	return r
}

func (svc *Service) privateEndpoints(ctx context.Context) *chi.Mux {
	log.Info(ctx, "enabling private endpoints")
	
	r := chi.NewRouter()

	// Middleware
	svc.identityClient = identity.New(svc.config.ZebedeeURL)
	checkIdentity := dphandlers.IdentityWithHTTPClient(svc.identityClient)
	permissions := middleware.NewPermissions(svc.config.ZebedeeURL, svc.config.EnablePermissionsAuth)

	r.Use(permissions.Require(auth.Permissions{Read: true}))
	r.Use(checkIdentity)
	r.Use(middleware.LogIdentity())

	// Routes
	hello := handler.NewHello(svc.responder, svc.cantabularClient)
	r.Get("/hello", hello.Get)
	r.Post("/hello", permissions.RequireCreate(hello.Create))

	return r
}
