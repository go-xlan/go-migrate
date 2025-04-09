package example1x

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-xlan/go-migrate/cobramigration"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

func TestNewMigrate(t *testing.T) {
	migrate := rese.P1(cobramigration.NewMigrate[*sqlite3.Sqlite](&cobramigration.Param{
		ScriptsInRoot: runpath.PARENT.Join("scripts"),
		ConnectSource: "sqlite3://file::memory:?cache=private",
	}))

	migrate.Log = &debugLogger{}

	require.NoError(t, migrate.Steps(+1))
	require.NoError(t, migrate.Steps(+1))
	require.NoError(t, migrate.Steps(-1))
	require.NoError(t, migrate.Steps(-1))
}

type debugLogger struct{}

func (l *debugLogger) Printf(format string, v ...interface{}) {
	fmt.Println(eroticgo.PINK.Sprint("->"), eroticgo.BLUE.Sprint(strings.TrimSpace(fmt.Sprintf(format, v...))))
}

func (l *debugLogger) Verbose() bool {
	return true // 启用详细日志
}
