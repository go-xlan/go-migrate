package example1_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-xlan/go-migrate/newmigrate"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

func TestNewMigrate(t *testing.T) {
	migration := rese.P1(newmigrate.NewMigrate[*sqlite3.Sqlite](&newmigrate.Param{
		ScriptsInRoot: runpath.PARENT.Join("scripts"),
		ConnectSource: "sqlite3://file::memory:?cache=private",
	}))

	migration.Log = &debugLogger{}

	caseShowVersion(t, migration)

	require.NoError(t, migration.Steps(+1))
	caseShowVersion(t, migration)

	require.NoError(t, migration.Steps(+1))
	caseShowVersion(t, migration)

	require.NoError(t, migration.Steps(-1))
	caseShowVersion(t, migration)

	require.NoError(t, migration.Steps(-1))
	caseShowVersion(t, migration)
}

type debugLogger struct{}

func (l *debugLogger) Printf(format string, v ...interface{}) {
	fmt.Println(eroticgo.PINK.Sprint("->"), eroticgo.BLUE.Sprint(strings.TrimSpace(fmt.Sprintf(format, v...))))
}

func (l *debugLogger) Verbose() bool {
	return true // 启用详细日志
}

func caseShowVersion(t *testing.T, migration *migrate.Migrate) {
	version, dirtyState, err := migration.Version()
	if err != nil {
		require.ErrorIs(t, err, migrate.ErrNilVersion)
	} else {
		require.NoError(t, err)
	}
	require.False(t, dirtyState)
	t.Log("version:", version)
}
