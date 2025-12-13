package migrationparam_test

import (
	"testing"

	"github.com/go-xlan/go-migrate/migrationparam"
)

func TestGetDebugMode(t *testing.T) {
	t.Log(migrationparam.GetDebugMode())
}
