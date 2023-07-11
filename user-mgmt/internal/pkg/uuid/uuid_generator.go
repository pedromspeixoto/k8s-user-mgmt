package uuid

import (
	"github.com/docker/distribution/uuid"
)

func GenerateUUID() string {
	return uuid.Generate().String()
}
