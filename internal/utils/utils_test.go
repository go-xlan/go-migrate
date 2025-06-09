package utils_test

import (
	"testing"

	"github.com/go-xlan/go-migrate/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestNewUUID32s(t *testing.T) {
	res := utils.NewUUID32s()
	t.Log(res)
	require.Len(t, res, 32)
}
