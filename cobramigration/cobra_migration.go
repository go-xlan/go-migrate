package cobramigration

import (
	"github.com/go-xlan/go-migrate/cobramigration/internal/cobramigrate"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
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

func NewMigrate(param *Param) (*migrate.Migrate, error) {
	migration, err := migrate.New("file://"+param.ScriptsInRoot, param.ConnectSource)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}

func NewMigrateCmd(migration *migrate.Migrate) *cobra.Command {
	return cobramigrate.NewMigrateCmd(migration)
}
