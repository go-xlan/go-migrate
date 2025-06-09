package sketch2

import (
	"embed"
	"testing"

	"github.com/go-xlan/go-migrate/internal/tests"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/rese"
)

//go:embed scripts/*.sql
var migrationsFS embed.FS

func TestEmbedFileSystem(t *testing.T) {
	// cp from https://github.com/golang-migrate/migrate/blob/278833935c12dda022b1355f33a897d895501c45/source/iofs/iofs_test.go#L14
	sourceInstance, err := iofs.New(migrationsFS, "scripts")
	require.NoError(t, err)
	defer rese.F0(sourceInstance.Close)

	version, err := sourceInstance.First()
	require.NoError(t, err)
	t.Log(version)

	tests.ShowSourceContent(t, sourceInstance, version, source.Up)
	tests.ShowSourceContent(t, sourceInstance, version, source.Down)
}

//go:embed scripts/00002_create_table_tb2.up.sql scripts/00002_create_table_tb2.down.sql
var migrationsFS2 embed.FS

func TestEmbedFileSystem2(t *testing.T) {
	sourceInstance, err := iofs.New(migrationsFS2, "scripts")
	require.NoError(t, err)
	defer rese.F0(sourceInstance.Close)

	version, err := sourceInstance.First()
	require.NoError(t, err)
	t.Log(version)

	tests.ShowSourceContent(t, sourceInstance, version, source.Up)
	tests.ShowSourceContent(t, sourceInstance, version, source.Down)
}
