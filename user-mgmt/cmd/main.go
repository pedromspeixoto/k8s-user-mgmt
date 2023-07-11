package main

import (
	"flag"
	"github.com/pedromspeixoto/users-api/internal/config"
	"github.com/pedromspeixoto/users-api/internal/data"
	"github.com/pedromspeixoto/users-api/internal/data/models"
	"github.com/pedromspeixoto/users-api/internal/domain"
	"github.com/pedromspeixoto/users-api/internal/http"
	"github.com/pedromspeixoto/users-api/internal/http/handlers"
	"github.com/pedromspeixoto/users-api/internal/pkg/files"
	"github.com/pedromspeixoto/users-api/internal/pkg/logger"
	"github.com/pedromspeixoto/users-api/internal/pkg/validator"
	"go.uber.org/fx"
)

// @title Users API
// @version 1.0
// @description Users API - Manage user and files
// @BasePath /
func main() {
	var cfgFilePath string
	flag.StringVar(
		&cfgFilePath,
		"config",
		"",
		"Path to config file. If not provided, config will be parsed from the environment.",
	)
	flag.Parse()

	app := fx.New(
		// Provide
		config.ProvideConfig(cfgFilePath),
		logger.ProvideLogger(),
		validator.ProvideValidator(),
		data.ProvideData(),
		models.ProvideModels(),
		domain.ProvideDomains(),
		handlers.ProvideHandlers(),
		files.ProvideFileServingClient(),
		// Invoke
		http.InvokeServer(),
	)

	app.Run()
}
