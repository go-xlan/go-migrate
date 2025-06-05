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
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewMigrate(t *testing.T) {
	db := rese.P1(gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	caseShowTableCount(t, db)
	caseShowTableCount(t, db)

	t.Log("new-migrate")

	migration := rese.P1(newmigrate.NewMigrate[*sqlite3.Sqlite](&newmigrate.Param{
		ScriptsInRoot: runpath.PARENT.Join("scripts"),
		ConnectSource: "sqlite3://file::memory:?cache=shared",
	}))
	defer func() {
		err1, err2 := migration.Close()
		must.Done(err1)
		must.Done(err2)
	}()

	migration.Log = &debugLogger{}

	t.Log("new-migrate-success")

	caseShowTableCount(t, db)
	caseShowVersionNum(t, migration, db)
	requireNotTable(t, db, "users")
	requireNotTable(t, db, "posts")

	t.Log("run-migrate")

	require.NoError(t, migration.Steps(+1))
	caseShowVersionNum(t, migration, db)
	requireHasTable(t, db, "users")
	requireNotTable(t, db, "posts")

	require.NoError(t, migration.Steps(+1))
	caseShowVersionNum(t, migration, db)
	requireHasTable(t, db, "users")
	requireHasTable(t, db, "posts")

	require.NoError(t, migration.Steps(-1))
	caseShowVersionNum(t, migration, db)
	requireHasTable(t, db, "users")
	requireNotTable(t, db, "posts")

	require.NoError(t, migration.Steps(-1))
	caseShowVersionNum(t, migration, db)
	requireNotTable(t, db, "users")
	requireNotTable(t, db, "posts")

	t.Log("run-migrate-success")
}

type debugLogger struct{}

func (l *debugLogger) Printf(format string, v ...interface{}) {
	fmt.Println(eroticgo.PINK.Sprint("->"), eroticgo.BLUE.Sprint(strings.TrimSpace(fmt.Sprintf(format, v...))))
}

func (l *debugLogger) Verbose() bool {
	return true // 启用详细日志
}

func caseShowVersionNum(t *testing.T, migration *migrate.Migrate, db *gorm.DB) {
	t.Log("---")
	version, dirtyState, err := migration.Version()
	if err != nil {
		require.ErrorIs(t, err, migrate.ErrNilVersion)
	} else {
		require.NoError(t, err)
	}
	require.False(t, dirtyState)
	t.Log("version-num:", version)

	caseShowTableCount(t, db)
	t.Log("---")
}

func caseShowTableCount(t *testing.T, db *gorm.DB) {
	tableList, err := db.Migrator().GetTables()
	require.NoError(t, err)
	t.Log("table-count:", len(tableList), tableList)
}

func requireHasTable(t *testing.T, db *gorm.DB, tableName string) {
	require.True(t, db.Migrator().HasTable(tableName))
}

func requireNotTable(t *testing.T, db *gorm.DB, tableName string) {
	require.False(t, db.Migrator().HasTable(tableName))
}
