package utils

import (
	"encoding/hex"

	"github.com/google/uuid"
)

func NewUUID32s() string {
	u := uuid.New()
	return hex.EncodeToString(u[:])
}
