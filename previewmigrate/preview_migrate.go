// Package previewmigrate: Migration preview functionality for safe SQL testing
// Provides zero-cost error recovery by testing migrations in transactions that always rollback
// Features intelligent script discovery and integration with existing migration workflows
// Prevents dirty state issues in automated migration processes
//
// previewmigrate: 用于安全 SQL 测试的迁移预览功能
// 通过在始终回滚的事务中测试迁移提供零成本错误恢复
// 具有智能脚本发现功能，与现有迁移工作流集成
// 防止自动化迁移过程中的脏状态问题
package previewmigrate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xlan/go-migrate/internal/utils"
	"github.com/go-xlan/go-migrate/newscripts"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"gorm.io/gorm"
)

// NewPreviewCmd creates preview command for migration dry-run with subcommands
// Uses lazy initialization - connections created only when command runs (not during command tree building)
// Migration factory accepts database connection parameter to share single connection (avoiding duplicate connections)
//
// NewPreviewCmd 创建具有子命令的迁移试运行预览命令
// 使用延迟初始化 - 仅在命令运行时创建连接（而非命令树构建时）
// 迁移工厂接受数据库连接参数以共享单个连接（避免重复连接）
func NewPreviewCmd(getDB func() *gorm.DB, getMigration func(*gorm.DB) *migrate.Migrate, scriptsPath string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "preview",
		Short: "Preview migrations (dry-run)",
		Long:  "Test migration SQL without applying changes",
	}

	rootCmd.AddCommand(newPreviewIncCmd(getDB, getMigration, scriptsPath))
	return rootCmd
}

// newPreviewIncCmd creates command for previewing next migration step
// Tests next migration SQL in transaction without applying changes to database
//
// newPreviewIncCmd 创建用于预览下一个迁移步骤的命令
// 在事务中测试下一个迁移 SQL 而不对数据库应用更改
func newPreviewIncCmd(getDB func() *gorm.DB, getMigration func(*gorm.DB) *migrate.Migrate, scriptsPath string) *cobra.Command {
	return &cobra.Command{
		Use:   "inc",
		Short: "Preview next migration step (+1)",
		Long:  "Test next migration SQL without applying changes",
		Run: func(cmd *cobra.Command, args []string) {
			db := getDB()
			migration := getMigration(db)
			err := previewNextMigration(migration, db, scriptsPath)
			if err != nil {
				zaplog.SUG.Debugln(eroticgo.RED.Sprint("PREVIEW FAILED:"))
				zaplog.SUG.Errorln(err)
				return
			}
		},
	}
}

// previewNextMigration previews the next migration without applying it
// Uses existing GetNextScriptInfo to find next script and tests execution in rollback transaction
// Provides comprehensive feedback on SQL validity and execution safety
//
// previewNextMigration 预览下一个迁移而不应用它
// 使用现有的 GetNextScriptInfo 找到下一个脚本并在回滚事务中测试执行
// 提供关于 SQL 有效性和执行安全性的全面反馈
func previewNextMigration(migration *migrate.Migrate, db *gorm.DB, scriptsPath string) error {
	// 1. Get current version
	currentVersion, dirtyFlag, err := migration.Version()
	utils.WhistleCause(err) //panic when cause is not expected
	if dirtyFlag {
		return erero.Errorf("DATABASE IS DIRTY AT VERSION %d", currentVersion)
	}

	// 2. Use existing GetNextScriptInfo to find next script
	options := newscripts.NewOptions(scriptsPath)
	scriptNaming := newscripts.NewScriptNaming()
	scriptInfo := newscripts.GetNextScriptInfo(migration, options, scriptNaming)

	scriptNames := scriptInfo.GetScriptNames()

	// Read the up script content
	forwardScriptPath := osmustexist.FILE(filepath.Join(scriptsPath, scriptNames.ForwardName))
	sqlContent := rese.V1(os.ReadFile(forwardScriptPath))
	if len(strings.TrimSpace(string(sqlContent))) == 0 {
		zaplog.SUG.Infoln(eroticgo.BLUE.Sprint("EMPTY MIGRATION FILE - PREVIEW SUCCESS"))
		return nil
	}

	zaplog.SUG.Infof("PREVIEWING MIGRATION SCRIPT: %s", scriptNames.ForwardName)

	// 3. Preview in transaction (always rollback)
	tx := db.Begin()
	if tx.Error != nil {
		return erero.Errorf("FAILED TO BEGIN TRANSACTION: %v", tx.Error)
	}

	// Execute and always rollback
	err = tx.Exec(string(sqlContent)).Error
	tx.Rollback() // Always rollback - this is a preview!

	if err != nil {
		zaplog.SUG.Debugln(eroticgo.RED.Sprint("PREVIEW FAILED - SQL EXECUTION ERROR:"))
		zaplog.SUG.Errorln(err)
		return erero.Errorf("PREVIEW FAILED: %v", err)
	}

	zaplog.SUG.Infoln(eroticgo.GREEN.Sprint("PREVIEW SUCCESS"))
	return nil
}
