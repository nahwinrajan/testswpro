package uuidgen

import (
	"github.com/google/uuid"
)

func NewRandom() (uuid.UUID, error) {
	uuid, err := uuid.NewRandom()

	return uuid, err
}
