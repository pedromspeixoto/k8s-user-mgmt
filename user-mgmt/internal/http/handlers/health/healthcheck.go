package health

import (
	"context"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi"
	health_domain "github.com/pedromspeixoto/users-api/internal/domain/health"
	"go.uber.org/fx"
)

type HealthServiceHandler interface {
	// Routes creates a REST router for the health service
	Routes() chi.Router
}

type healthServiceDeps struct {
	fx.In

	HealthService health_domain.HealthService
}

type healthServiceHandler struct {
	healthServiceDeps
}

func NewHealthServiceHandler(deps healthServiceDeps) HealthServiceHandler {
	return &healthServiceHandler{
		healthServiceDeps: deps,
	}
}

func (h healthServiceHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", health.NewHandler(GetChecker(h.healthServiceDeps)))

	return r
}

// GetChecker - Check service status
// @Summary Get service status.
// @Description This API is used to get the environment and dependencies status.
// @Tags health
// @Accept  json
// @Produce  json
// @Router /health [get]
func GetChecker(deps healthServiceDeps) health.Checker {
	// Create a new Checker.
	return health.NewChecker(

		// Set the time-to-live for our cache to 1 second (default).
		health.WithCacheDuration(1*time.Second),

		// Configure a global timeout that will be applied to all checks.
		health.WithTimeout(10*time.Second),

		// A check configuration to see if our database connection is up.
		// The check function will be executed for each HTTP request.
		health.WithCheck(health.Check{
			Name:    "database",      // A unique check name.
			Timeout: 2 * time.Second, // A check specific timeout.
			Check: func(ctx context.Context) error {
				dbStatus := deps.HealthService.GetDbStatus()
				if dbStatus != nil {
					return dbStatus
				}
				return nil
			},
		}),
	)
}
