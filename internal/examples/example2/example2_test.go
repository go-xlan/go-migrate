package example2

import (
	"fmt"
	"testing"

	"github.com/go-xlan/go-migrate/checkmigration"
	"github.com/go-xlan/go-migrate/internal/tests"
	"github.com/go-xlan/go-migrate/newmigrate"
	"github.com/go-xlan/go-migrate/newscripts"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGenerateNewScript(t *testing.T) {
	scriptsInRoot := runpath.PARENT.Join("scripts/case1")

	db := rese.P1(gorm.Open(sqlite.Open(newDsn()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	migration := rese.P1(newmigrate.NewWithScriptsAndDatabase(
		&newmigrate.ScriptsAndDatabaseParam{
			ScriptsInRoot:    scriptsInRoot,
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

	options := newscripts.NewOptions(scriptsInRoot)
	require.False(t, options.DryRun)
	require.False(t, options.SurveyWritten)
	require.Equal(t, scriptsInRoot, options.ScriptsInRoot)

	scriptInfo := newscripts.GetNextScriptInfo(migration, options, newscripts.NewScriptNaming())
	t.Log(neatjsons.S(scriptInfo))
	require.Equal(t, newscripts.UpdateScript, scriptInfo.Action)
	require.Equal(t, "00001_script.up.sql", scriptInfo.ForwardName)
	require.Equal(t, "00001_script.down.sql", scriptInfo.ReverseName)

	migrateOps := checkmigration.GetMigrateOps(db, []any{
		&UserV1{},
	})
	scriptInfo.WriteScripts(migrateOps, options)
}

func newDsn() string {
	// 通常建议使用 shared (线程间共享/连接间共享)模式的，但使用 shared 模式时也通常建议设置唯一名称，就能控制共享的范围
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	return dsn
}

func TestGenerateNewScript_2(t *testing.T) {
	scriptsInRoot := runpath.PARENT.Join("scripts/case2")

	db := rese.P1(gorm.Open(sqlite.Open(newDsn()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	migration := rese.P1(newmigrate.NewWithScriptsAndDatabase(
		&newmigrate.ScriptsAndDatabaseParam{
			ScriptsInRoot:    scriptsInRoot,
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
	must.Done(migration.Steps(+1))

	require.True(t, t.Run("update-00002", func(t *testing.T) {
		options := newscripts.NewOptions(scriptsInRoot)
		options.DryRun = true
		scriptInfo := newscripts.GetNextScriptInfo(migration, options, newscripts.NewScriptNaming())
		require.Equal(t, newscripts.UpdateScript, scriptInfo.Action)
		require.Equal(t, "00002_script.up.sql", scriptInfo.ForwardName)
		require.Equal(t, "00002_script.down.sql", scriptInfo.ReverseName)

		migrateOps := checkmigration.GetMigrateOps(db, []any{
			&UserV2{},
		})
		scriptInfo.WriteScripts(migrateOps, options)
	}))

	must.Done(migration.Steps(+1))

	require.True(t, t.Run("create-00003", func(t *testing.T) {
		options := newscripts.NewOptions(scriptsInRoot)
		options.DryRun = true
		scriptInfo := newscripts.GetNextScriptInfo(migration, options, newscripts.NewScriptNaming())
		require.Equal(t, newscripts.CreateScript, scriptInfo.Action)
		require.Equal(t, "00003_script.up.sql", scriptInfo.ForwardName)
		require.Equal(t, "00003_script.down.sql", scriptInfo.ReverseName)

		migrationOps := checkmigration.GetMigrateOps(db, []any{
			&UserV2{},
		})
		scriptInfo.WriteScripts(migrationOps, options)
	}))
}

func TestGenerateNewScript_3(t *testing.T) {
	scriptsInRoot := runpath.PARENT.Join("scripts/case3")

	db := rese.P1(gorm.Open(sqlite.Open(newDsn()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	migration := rese.P1(newmigrate.NewWithScriptsAndDatabase(
		&newmigrate.ScriptsAndDatabaseParam{
			ScriptsInRoot:    scriptsInRoot,
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
	must.Done(migration.Steps(+1))

	options := newscripts.NewOptions(scriptsInRoot)
	scriptInfo := newscripts.GetNextScriptInfo(migration, options, newscripts.NewScriptNaming())
	require.Equal(t, newscripts.UpdateScript, scriptInfo.Action)
	require.Equal(t, "00002_script.up.sql", scriptInfo.ForwardName)
	require.Equal(t, "00002_script.down.sql", scriptInfo.ReverseName)

	migrationOps := checkmigration.GetMigrateOps(db, []any{
		&UserV2{},
	})
	scriptInfo.WriteScripts(migrationOps, options)
}
