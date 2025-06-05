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

type ScriptsAndSourceParam struct {
	ScriptsInRoot string
	ConnectSource string
}

// NewWithScriptsAndSource creates a new migration instance.
// T must implement database.Driver interface, e.g.:
//   - sqlite3.Sqlite from "github.com/golang-migrate/migrate/v4/database/sqlite3"
//   - mysql.Mysql from "github.com/golang-migrate/migrate/v4/database/mysql"
//   - postgres.Postgres from "github.com/golang-migrate/migrate/v4/database/postgres"
//
// Note: The type param T is used to trigger side effects.
//
// Example:
//   - migration, err := NewWithScriptsAndSource[*sqlite3.Sqlite](param)
//   - migration, err := NewWithScriptsAndSource[*mysql.Mysql](param)
//   - migration, err := NewWithScriptsAndSource[*postgres.Postgres](param)
func NewWithScriptsAndSource[T database.Driver](param *ScriptsAndSourceParam) (*migrate.Migrate, error) {
	migration, err := migrate.New("file://"+param.ScriptsInRoot, param.ConnectSource)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}

type EmbedAndDatabaseParam struct {
	MigrationsFS     *embed.FS
	EmbedDirName     string
	DatabaseName     string
	DatabaseInstance database.Driver
}

func NewWithEmbedAndDatabase(param *EmbedAndDatabaseParam) (*migrate.Migrate, error) {
	migration, err := migrate.NewWithInstance(
		"iofs", // 固定的 iofs 类型
		rese.V1(iofs.New(param.MigrationsFS, param.EmbedDirName)), // 初始化 iofs 驱动
		param.DatabaseName,
		param.DatabaseInstance,
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}
