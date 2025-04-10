package newmigrate

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
)

func init() {
	must.Full(&file.File{})
}

type Param struct {
	ScriptsInRoot string
	ConnectSource string
}

// NewMigrate creates a new migration instance.
// T must implement database.Driver interface, e.g.:
//   - sqlite3.Sqlite from "github.com/golang-migrate/migrate/v4/database/sqlite3"
//   - mysql.Mysql from "github.com/golang-migrate/migrate/v4/database/mysql"
//   - postgres.Postgres from "github.com/golang-migrate/migrate/v4/database/postgres"
//
// Note: The type param T is used to trigger side effects.
//
// Example:
//   - migration, err := NewMigrate[*sqlite3.Sqlite](param)
//   - migration, err := NewMigrate[*mysql.Mysql](param)
//   - migration, err := NewMigrate[*postgres.Postgres](param)
func NewMigrate[T database.Driver](param *Param) (*migrate.Migrate, error) {
	migration, err := migrate.New("file://"+param.ScriptsInRoot, param.ConnectSource)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}
