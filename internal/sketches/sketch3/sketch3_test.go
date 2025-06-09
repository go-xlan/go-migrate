package sketch3

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
)

//go:embed scripts
var migrationsFS embed.FS

func TestEmbedFileSystem(t *testing.T) {
	entries, err := fs.ReadDir(migrationsFS, "scripts")
	require.NoError(t, err)

	for _, e := range entries {
		require.False(t, e.IsDir())

		migration, err := source.DefaultParse(e.Name())
		require.NoError(t, err)
		t.Log(neatjsons.S(migration))
	}
}

//go:embed scripts/00002_create_table_tb2.up.sql scripts/00002_create_table_tb2.down.sql
var migrationsFS2 embed.FS

func TestEmbedFileSystem2(t *testing.T) {
	entries, err := fs.ReadDir(migrationsFS2, "scripts")
	require.NoError(t, err)

	// cp from https://github.com/golang-migrate/migrate/blob/278833935c12dda022b1355f33a897d895501c45/source/iofs/iofs.go#L55
	ms := source.NewMigrations()
	for _, e := range entries {
		require.False(t, e.IsDir())

		migration, err := source.DefaultParse(e.Name())
		require.NoError(t, err)
		t.Log("append migration to migrations:", "version:", migration.Version, "direction:", migration.Direction)
		require.True(t, ms.Append(migration))
	}

	{
		version, ok := ms.First()
		require.True(t, ok)
		t.Log(version)
		require.Equal(t, uint64(2), uint64(version))
	}
	{
		const version = 2
		migration, ok := ms.Up(version)
		require.True(t, ok)
		require.NotNil(t, migration)
		t.Log("show migrate up version:", version, neatjsons.S(migration))
		require.Equal(t, uint64(version), uint64(migration.Version))
	}
	{
		const version = 2
		migration, ok := ms.Down(version)
		require.True(t, ok)
		require.NotNil(t, migration)
		t.Log("show migrate down version:", version, neatjsons.S(migration))
		require.Equal(t, uint64(version), uint64(migration.Version))
	}
}
