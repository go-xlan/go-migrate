package cobrasqlite3migration

import (
	"github.com/go-xlan/go-migrate/cobramigration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
)

func init() {
	must.Full(&sqlite3.Sqlite{})
}

func New(param *cobramigration.Param) (*migrate.Migrate, error) {
	migration, err := cobramigration.NewMigrate(param)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}

func NewMigrateCmd(migration *migrate.Migrate) *cobra.Command {
	return cobramigration.NewMigrateCmd(migration)
}
