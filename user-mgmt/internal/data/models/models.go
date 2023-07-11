package models

import (
	"github.com/pedromspeixoto/users-api/internal/data/models/users"
	"go.uber.org/fx"
)

func ProvideModels() fx.Option {
	return fx.Options(
		fx.Provide(
			users.NewUserRepository,
			users.NewUserFileRepository,
		),
	)
}
