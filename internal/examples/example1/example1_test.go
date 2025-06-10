package example1_test

import (
	"embed"
	"fmt"
	"testing"

	"github.com/go-xlan/go-migrate/internal/tests"
	"github.com/go-xlan/go-migrate/newmigrate"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewWithScriptsAndDBSource(t *testing.T) {
	migration := rese.P1(newmigrate.NewWithScriptsAndDBSource[*sqlite3.Sqlite](&newmigrate.ScriptsAndDBSourceParam{
		ScriptsInRoot: runpath.PARENT.Join("scripts"),
		ConnectSource: "sqlite3://file::memory:?cache=private",
	}))
	migration.Log = &tests.LoggerDebug{}
	defer func() {
		err1, err2 := migration.Close()
		must.Done(err1)
		must.Done(err2)
	}()
	must.Done(migration.Up())

	tests.CaseShowVersionNum(t, migration)
}

func TestNewWithScriptsAndDatabase(t *testing.T) {
	// 通常建议使用 shared (线程间共享/连接间共享)模式的，但使用 shared 模式时也通常建议设置唯一名称，就能控制共享的范围
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	tests.CaseShowTableCount(t, db)

	t.Log("new-migrate")

	migration := rese.P1(newmigrate.NewWithScriptsAndDatabase(
		&newmigrate.ScriptsAndDatabaseParam{
			ScriptsInRoot:    runpath.PARENT.Join("scripts"),
			DatabaseName:     "sqlite3",
			DatabaseInstance: rese.V1(sqlite3.WithInstance(rese.P1(db.DB()), &sqlite3.Config{})),
		},
	))
	migration.Log = &tests.LoggerDebug{}
	defer func() {
		err1, err2 := migration.Close()
		must.Done(err1)
		must.Done(err2)
	}()

	t.Log("new-migrate-success")

	tests.CaseShowTableCount(t, db)
	tests.CaseShowVersionNumAndTables(t, migration, db)
	tests.RequireNotTable(t, db, "users")
	tests.RequireNotTable(t, db, "posts")

	t.Log("run-migrate")

	require.NoError(t, migration.Steps(+1))
	tests.CaseShowVersionNumAndTables(t, migration, db)
	tests.RequireHasTable(t, db, "users")
	tests.RequireNotTable(t, db, "posts")

	require.NoError(t, migration.Steps(+1))
	tests.CaseShowVersionNumAndTables(t, migration, db)
	tests.RequireHasTable(t, db, "users")
	tests.RequireHasTable(t, db, "posts")

	require.NoError(t, migration.Steps(-1))
	tests.CaseShowVersionNumAndTables(t, migration, db)
	tests.RequireHasTable(t, db, "users")
	tests.RequireNotTable(t, db, "posts")

	require.NoError(t, migration.Steps(-1))
	tests.CaseShowVersionNumAndTables(t, migration, db)
	tests.RequireNotTable(t, db, "users")
	tests.RequireNotTable(t, db, "posts")

	t.Log("run-migrate-success")
}

//go:embed scripts
var migrationsFS embed.FS // 使用这个也行 go:embed scripts/*.sql 而且更精确些，但通常认为这个目录就只有 sql 类型文件，没有别的

func TestNewWithEmbedFsAndDatabase(t *testing.T) {
	// 数据库连接字符串，使用内存中的 SQLite 数据库
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	// 创建迁移实例
	migration := rese.P1(newmigrate.NewWithEmbedFsAndDatabase(
		&newmigrate.EmbedFsAndDatabaseParam{
			MigrationsFS:     &migrationsFS,
			EmbedDirName:     "scripts",
			DatabaseName:     "sqlite3",
			DatabaseInstance: rese.V1(sqlite3.WithInstance(rese.P1(db.DB()), &sqlite3.Config{})),
		},
	))
	migration.Log = &tests.LoggerDebug{}
	defer func() {
		err1, err2 := migration.Close()
		must.Done(err1)
		must.Done(err2)
	}()

	// 执行迁移
	require.NoError(t, migration.Up())
	t.Log("success")

	tests.RequireHasTables(t, db, []string{"users", "posts"})
}
