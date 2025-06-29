package tests

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoggerDebug struct{}

func (l *LoggerDebug) Printf(format string, values ...interface{}) {
	fmt.Println(eroticgo.PINK.Sprint("->"), eroticgo.BLUE.Sprint(strings.TrimSpace(fmt.Sprintf(format, values...))))
}

func (l *LoggerDebug) Verbose() bool {
	return true // 启用详细日志
}

func CaseShowVersionNumAndTables(t *testing.T, migration *migrate.Migrate, db *gorm.DB) {
	zapLog := zaplog.ZAPS.Skip(1)

	version := caseShowVersionNum(t, migration, zapLog.SkipZap(1))

	caseShowTableCount(t, db, zapLog.SkipZap(1, zap.Uint("version", version)))
	t.Log("---")
}

func CaseShowVersionNum(t *testing.T, migration *migrate.Migrate) {
	zapLog := zaplog.ZAPS.Skip(1)
	caseShowVersionNum(t, migration, zapLog.SkipZap(1))
	t.Log("---")
}

func caseShowVersionNum(t *testing.T, migration *migrate.Migrate, zapLog *zaplog.Zap) uint {
	t.Log("---")
	version, dirtyFlag, err := migration.Version()
	if err != nil {
		require.ErrorIs(t, err, migrate.ErrNilVersion)
	} else {
		require.NoError(t, err)
	}
	require.False(t, dirtyFlag)
	zapLog.SUG.Debugln("version-num:", version)
	return version
}

func CaseShowTableCount(t *testing.T, db *gorm.DB) {
	zapLog := zaplog.ZAPS.Skip(1)
	caseShowTableCount(t, db, zapLog.SkipZap(1))
	t.Log("---")
}

func caseShowTableCount(t *testing.T, db *gorm.DB, zapLog *zaplog.Zap) {
	tableList, err := db.Migrator().GetTables()
	require.NoError(t, err)
	zapLog.SUG.Debugln("table-count:", len(tableList), tableList)
}

func RequireHasTable(t *testing.T, db *gorm.DB, tableName string) {
	require.True(t, db.Migrator().HasTable(tableName))
}

func RequireNotTable(t *testing.T, db *gorm.DB, tableName string) {
	require.False(t, db.Migrator().HasTable(tableName))
}

func RequireHasTables(t *testing.T, db *gorm.DB, tableNames []string) {
	for _, tableName := range tableNames {
		RequireHasTable(t, db, tableName)
	}
}

func ShowSourceContent(t *testing.T, sourceInstance source.Driver, version uint, direction source.Direction) {
	readFunc := tern.BVV(direction == source.Up, sourceInstance.ReadUp, sourceInstance.ReadDown)
	scriptFile, name, err := readFunc(version)
	require.NoError(t, err)
	defer rese.F0(scriptFile.Close)
	t.Log(name, direction)
	content, err := io.ReadAll(scriptFile)
	require.NoError(t, err)
	t.Log(string(content))
}
