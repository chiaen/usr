package uuid

import (
	"github.com/google/uuid"
)

func New() (uuid.UUID, error) {
	return uuid.NewRandom()
}
