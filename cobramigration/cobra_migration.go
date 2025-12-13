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
	"github.com/go-xlan/go-migrate/migrationparam"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
)

// NewMigrateCmd creates comprehensive migration command with subcommands for all migration operations
// Uses lazy initialization - connections created only when command runs (not during command tree building)
// Migration connection interface ensures proper resource cleanup after operations
//
// NewMigrateCmd 创建包含子命令的综合迁移命令，用于所有迁移操作
// 使用延迟初始化 - 仅在命令运行时创建连接（而非命令树构建时）
// 迁移连接接口确保操作后正确清理资源
func NewMigrateCmd(param *migrationparam.MigrationParam) *cobra.Command {
	// Create root command
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Database migration",
		Long:  "Database migration",
		Run: func(cmd *cobra.Command, args []string) {
			migration, cleanup := param.GetMigration()
			defer cleanup()

			version, dirtyFlag, err := migration.Version()
			utils.WhistleCause(err) //panic when cause is not expected
			if dirtyFlag {
				eroticgo.RED.ShowMessage(version, "(DIRTY)")
			} else {
				eroticgo.GREEN.ShowMessage(version)
			}
		},
	}

	rootCmd.AddCommand(newAllCmd(param)) // Add `all` command
	rootCmd.AddCommand(newIncCMD(param)) // Add `inc` command
	rootCmd.AddCommand(newDecCMD(param)) // Add `dec` command

	return rootCmd
}

// newAllCmd creates command for executing all pending migrations
// Performs complete database upgrade to latest schema version
//
// newAllCmd 创建用于执行所有待处理迁移的命令
// 将数据库升级到最新的结构版本
func newAllCmd(param *migrationparam.MigrationParam) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Run all migration files",
		Run: func(cmd *cobra.Command, args []string) {
			migration, cleanup := param.GetMigration()
			defer cleanup()

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
func newDecCMD(param *migrationparam.MigrationParam) *cobra.Command {
	return &cobra.Command{
		Use:   "dec",
		Short: "Rollback one step (-1)",
		Run: func(cmd *cobra.Command, args []string) {
			migration, cleanup := param.GetMigration()
			defer cleanup()

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
func newIncCMD(param *migrationparam.MigrationParam) *cobra.Command {
	return &cobra.Command{
		Use:   "inc",
		Short: "Run next step (+1)",
		Run: func(cmd *cobra.Command, args []string) {
			migration, cleanup := param.GetMigration()
			defer cleanup()

			// Execute next migration step forward
			// 向前执行下一个迁移步骤
			utils.WhistleCause(migration.Steps(+1))
		},
	}
}
