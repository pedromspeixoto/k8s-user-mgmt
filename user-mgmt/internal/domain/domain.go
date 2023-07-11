package domain

import (
	"go.uber.org/fx"

	"github.com/pedromspeixoto/users-api/internal/domain/health"
	"github.com/pedromspeixoto/users-api/internal/domain/users"
)

func ProvideDomains() fx.Option {
	return fx.Provide(
		health.NewHealthService,
		users.NewUserService,
	)
}
