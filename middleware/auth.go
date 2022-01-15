package middleware

import (
	"net/http"

	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/log.go/v2/log"
)

func IsAuthenticated() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			ctx := r.Context()

			if id := request.Caller(ctx); id != ""{
				log.Info(ctx, "caller identity verified", log.Data{
					"caller_identity": id,
					"URI": r.URL.Path,
				})
				next.ServeHTTP(w, r)
				return
			}

			log.Info(ctx, "caller identity not present", log.Data{
				"URI": r.URL.Path,
			})
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		})
	}
}

type AuthHandler interface {
	Require(auth.Permissions, http.HandlerFunc) http.HandlerFunc
}

type Permissions struct {
	handler AuthHandler
}

func NewPermissions(zebedeeURL string, enabled bool) *Permissions {
	if !enabled{
		return &Permissions{
			handler: &auth.NopHandler{},
		}
	}

	client := auth.NewPermissionsClient(dphttp.NewClient())
	verifier := auth.DefaultPermissionsVerifier()

	return &Permissions{
		handler: auth.NewHandler(
			auth.NewPermissionsRequestBuilder(zebedeeURL),
			client,
			verifier,
		),
	}
}

func (p *Permissions) Require(required auth.Permissions) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			p.handler.Require(required, next.ServeHTTP)(w, r)
		})
	}
}
