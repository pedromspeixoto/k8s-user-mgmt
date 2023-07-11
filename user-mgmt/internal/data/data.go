package data

import (
	"go.uber.org/fx"
)

func ProvideData() fx.Option {
	return fx.Provide(
		NewDbClient,
	)
}
