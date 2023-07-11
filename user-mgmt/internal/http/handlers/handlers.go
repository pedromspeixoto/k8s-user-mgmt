package handlers

import (
	"github.com/pedromspeixoto/users-api/internal/http/handlers/health"
	"github.com/pedromspeixoto/users-api/internal/http/handlers/metrics"
	"github.com/pedromspeixoto/users-api/internal/http/handlers/users"
	"go.uber.org/fx"
)

func ProvideHandlers() fx.Option {
	return fx.Provide(
		health.NewHealthServiceHandler,
		metrics.NewMetricServiceHandler,
		users.NewUserServiceHandler,
	)
}
