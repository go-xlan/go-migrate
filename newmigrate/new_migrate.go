package newmigrate

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func init() {
	must.Full(&file.File{}) // register file source(side effects)
}

type ScriptsAndDBSourceParam struct {
	ScriptsInRoot string
	ConnectSource string
}

// NewWithScriptsAndDBSource creates a new migration instance.
// T must implement database.Driver interface, e.g.:
//   - sqlite3.Sqlite from "github.com/golang-migrate/migrate/v4/database/sqlite3"
//   - mysql.Mysql from "github.com/golang-migrate/migrate/v4/database/mysql"
//   - postgres.Postgres from "github.com/golang-migrate/migrate/v4/database/postgres"
//
// Note: The type param T is used to trigger side effects.
//
// Example:
//   - migration, err := NewWithScriptsAndDBSource[*sqlite3.Sqlite](param)
//   - migration, err := NewWithScriptsAndDBSource[*mysql.Mysql](param)
//   - migration, err := NewWithScriptsAndDBSource[*postgres.Postgres](param)
func NewWithScriptsAndDBSource[T database.Driver](param *ScriptsAndDBSourceParam) (*migrate.Migrate, error) {
	sourceURL := "file://" + param.ScriptsInRoot
	migration, err := migrate.New(
		sourceURL,
		param.ConnectSource,
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}

type ScriptsAndDatabaseParam struct {
	ScriptsInRoot    string
	DatabaseName     string
	DatabaseInstance database.Driver
}

func NewWithScriptsAndDatabase(param *ScriptsAndDatabaseParam) (*migrate.Migrate, error) {
	sourceURL := "file://" + param.ScriptsInRoot
	migration, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		param.DatabaseName,
		param.DatabaseInstance,
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}

type EmbedFsAndDatabaseParam struct {
	MigrationsFS     *embed.FS
	EmbedDirName     string
	DatabaseName     string
	DatabaseInstance database.Driver
}

func NewWithEmbedFsAndDatabase(param *EmbedFsAndDatabaseParam) (*migrate.Migrate, error) {
	const sourceName = "iofs"
	// cp from https://github.com/golang-migrate/migrate/blob/278833935c12dda022b1355f33a897d895501c45/source/iofs/example_test.go#L22
	migration, err := migrate.NewWithInstance(
		sourceName, // 固定的 iofs 类型
		rese.V1(iofs.New(param.MigrationsFS, param.EmbedDirName)), // 初始化 iofs 驱动
		param.DatabaseName,
		param.DatabaseInstance,
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}
