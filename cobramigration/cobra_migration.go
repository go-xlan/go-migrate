package cobramigration

import (
	"github.com/go-xlan/go-migrate/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
)

func NewMigrateCmd(migration *migrate.Migrate) *cobra.Command {
	// Create root command
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Database migration",
		Long:  "Database migration",
		Run: func(cmd *cobra.Command, args []string) {
			version, dirtyFlag, err := migration.Version()
			utils.WhistleCause(err) //panic when cause is not expected
			if dirtyFlag {
				eroticgo.RED.ShowMessage(version, "(DIRTY)")
			} else {
				eroticgo.GREEN.ShowMessage(version)
			}
		},
	}

	rootCmd.AddCommand(newAllCmd(migration)) // Add `all` command
	rootCmd.AddCommand(newIncCMD(migration)) // Add `inc` command
	rootCmd.AddCommand(newDecCMD(migration)) // Add `dec` command

	return rootCmd
}

func newAllCmd(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Run all migration files",
		Run: func(cmd *cobra.Command, args []string) {
			utils.WhistleCause(migration.Up()) // Perform upgrade
		},
	}
}

func newDecCMD(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "dec",
		Short: "Rollback one step (-1)",
		Run: func(cmd *cobra.Command, args []string) {
			utils.WhistleCause(migration.Steps(-1))
		},
	}
}

func newIncCMD(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "inc",
		Short: "Run next step (+1)",
		Run: func(cmd *cobra.Command, args []string) {
			utils.WhistleCause(migration.Steps(+1))
		},
	}
}
