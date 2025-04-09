package cobramigrate

import (
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/zaplog"
)

func NewMigrateCmd(migration *migrate.Migrate) *cobra.Command {
	// Create root command
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Database migration",
		Long:  "Database migration",
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
			whistle(migration.Up()) // Perform upgrade
		},
	}
}

func newDecCMD(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "dec",
		Short: "Rollback one step (-1)",
		Run: func(cmd *cobra.Command, args []string) {
			whistle(migration.Steps(-1))
		},
	}
}

func newIncCMD(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "inc",
		Short: "Run next step (+1)",
		Run: func(cmd *cobra.Command, args []string) {
			whistle(migration.Steps(+1))
		},
	}
}

func whistle(ein error) {
	if ein != nil {
		if errors.Is(ein, migrate.ErrNoChange) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("NO MIGRATION FILES TO RUN"))
		} else if errors.Is(ein, os.ErrNotExist) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("MIGRATION FILES NOT FOUND"))
		} else {
			zaplog.SUG.Panicln(eroticgo.RED.Sprint("MIGRATION FAILED:"), ein)
		}
		return
	}
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint("MIGRATION SUCCESS"))
}
