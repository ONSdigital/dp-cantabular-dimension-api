package service

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-dimension-api/handler"
)

func (svc *Service) buildRoutes(){
	r := chi.NewRouter()

	hello := handler.NewHello(svc.responder)
	r.Get("/hello", hello.Get)

	r.Handle("/health", http.HandlerFunc(svc.healthCheck.Handler))

	if svc.config.EnablePrivateEndpoints{
		r.Mount("/", func() chi.Router{
			r := chi.NewRouter()
			r.Post("/hello", hello.Post)
			return r
		}())
	}

	svc.router = r
}
