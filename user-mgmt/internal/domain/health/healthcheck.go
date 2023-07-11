package health

import (
	"errors"
	"fmt"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	ErrCreateDb = errors.New("error creating database connection")
	ErrPingDb   = errors.New("error pinging database")
)

// HealthService provides methods pertaining to managing environments.
type HealthService interface {
	// get db status
	GetDbStatus() error
}

type healthServiceDeps struct {
	fx.In

	Db *gorm.DB
}

type healthService struct {
	healthServiceDeps
}

func NewHealthService(deps healthServiceDeps) HealthService {
	return &healthService{
		healthServiceDeps: deps,
	}
}

func (hs *healthService) GetDbStatus() error {
	sqldb, err := hs.healthServiceDeps.Db.DB()
	if err != nil {
		return fmt.Errorf("error creating database connection")
	}
	if err = sqldb.Ping(); err != nil {
		return fmt.Errorf("error pinging database")
	}
	return nil
}
