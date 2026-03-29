package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/handlers"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api/middleware"
	domainauth "github.com/k3forx/CoffeeBeansStock/backend/internal/domain/auth"
)

// Deps holds the dependencies required to set up routes.
type Deps struct {
	AuthHandler         *handlers.AuthHandler
	CoffeeBeansHandler  *handlers.CoffeeBeansHandler
	UsageHistoryHandler *handlers.UsageHistoryHandler
	TokenManager        domainauth.TokenManager
	HealthCheck         http.HandlerFunc
}

// New creates and configures the application router.
func New(deps Deps) chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestLogger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", deps.HealthCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.AuthHandler.RegisterAnonymous)
			r.Post("/refresh", deps.AuthHandler.Refresh)

			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(deps.TokenManager))
				r.Get("/me", deps.AuthHandler.GetMe)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.TokenManager))
			r.Route("/coffee-beans", func(r chi.Router) {
				r.Get("/", deps.CoffeeBeansHandler.List)
				r.Post("/", deps.CoffeeBeansHandler.Create)
				r.Get("/{id}", deps.CoffeeBeansHandler.Get)
				r.Put("/{id}", deps.CoffeeBeansHandler.Update)
				r.Delete("/{id}", deps.CoffeeBeansHandler.Delete)

				r.Route("/{id}/usages", func(r chi.Router) {
					r.Post("/", deps.UsageHistoryHandler.Create)
					r.Get("/", deps.UsageHistoryHandler.List)
					r.Delete("/{usageId}", deps.UsageHistoryHandler.Delete)
				})
			})
		})
	})

	return r
}
