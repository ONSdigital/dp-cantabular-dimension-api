package service

import (
	"net/http"
	"context"

	"github.com/ONSdigital/dp-cantabular-dimension-api/handler"
	"github.com/ONSdigital/dp-cantabular-dimension-api/middleware"

	"github.com/ONSdigital/dp-authorisation/auth"
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
	r.Post("/hello", hello.Post)

	return r
}

func (svc *Service) privateEndpoints(ctx context.Context) *chi.Mux {
	log.Info(ctx, "enabling private endpoints")
	r := chi.NewRouter()

	permissions := middleware.NewPermissions(svc.config.ZebedeeURL, svc.config.EnablePermissionsAuth)
	// middleware called in FIFO order
	if svc.config.EnableIdentityAuth{
		r.Use(middleware.IsAuthenticated())
	}
	r.Use(permissions.Require(auth.Permissions{Read: true}))

	hello := handler.NewHello(svc.responder, svc.cantabularClient)
	r.Get("/hello", hello.Get)
	r.Post("/hello", hello.Post)

	return r
}
