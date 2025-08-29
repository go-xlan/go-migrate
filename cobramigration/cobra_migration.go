// Package cobramigration: Cobra CLI integration for database migration operations
// Provides command-line interface for migration execution with user-friendly commands
// Features version display, batch migration, and step-by-step migration control
// Integrates seamlessly with golang-migrate for robust migration management
//
// cobramigration: 用于数据库迁移操作的 Cobra CLI 集成
// 为迁移执行提供命令行接口，具有用户友好的命令
// 具有版本显示、批量迁移和逐步迁移控制功能
// 与 golang-migrate 无缝集成，提供稳健的迁移管理
package cobramigration

import (
	"github.com/go-xlan/go-migrate/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
)

// NewMigrateCmd creates comprehensive migration command with subcommands for all migration operations
// Provides root command that displays current migration version and dirty state
// Includes subcommands for batch migration, incremental steps, and rollback operations
//
// NewMigrateCmd 创建包含子命令的综合迁移命令，用于所有迁移操作
// 提供显示当前迁移版本和脏状态的根命令
// 包含用于批量迁移、增量步骤和回滚操作的子命令
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

// newAllCmd creates command for executing all pending migrations
// Performs complete database upgrade to latest schema version
//
// newAllCmd 创建用于执行所有待处理迁移的命令
// 将数据库升级到最新的结构版本
func newAllCmd(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Run all migration files",
		Run: func(cmd *cobra.Command, args []string) {
			// Perform complete database upgrade
			// 执行完整的数据库升级
			utils.WhistleCause(migration.Up())
		},
	}
}

// newDecCMD creates command for rolling back one migration step
// Safely reverts database schema by one version
//
// newDecCMD 创建用于回滚一个迁移步骤的命令
// 安全地将数据库结构回退一个版本
func newDecCMD(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "dec",
		Short: "Rollback one step (-1)",
		Run: func(cmd *cobra.Command, args []string) {
			// Rollback database by one migration step
			// 将数据库回滚一个迁移步骤
			utils.WhistleCause(migration.Steps(-1))
		},
	}
}

// newIncCMD creates command for executing next migration step
// Advances database schema by one version forward
//
// newIncCMD 创建用于执行下一个迁移步骤的命令
// 将数据库结构向前推进一个版本
func newIncCMD(migration *migrate.Migrate) *cobra.Command {
	return &cobra.Command{
		Use:   "inc",
		Short: "Run next step (+1)",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute next migration step forward
			// 向前执行下一个迁移步骤
			utils.WhistleCause(migration.Steps(+1))
		},
	}
}
