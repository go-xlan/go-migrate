// Package migrationstate: Database migration status inspection and reporting system
// Provides comprehensive status view of database version, script versions and schema differences
// Enables users to understand current migration state before performing operations
//
// migrationstate: 数据库迁移状态检查和报告系统
// 提供数据库版本、脚本版本和结构差异的综合状态视图
// 使用户在执行操作前了解当前迁移状态
package migrationstate

import (
	"fmt"
	"os"
	"sort"

	"github.com/go-xlan/go-migrate/checkmigration"
	"github.com/go-xlan/go-migrate/newmigrate"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/rese"
	"gorm.io/gorm"
)

// Config contains configuration options the status command needs
// Provides dependencies required when performing status inspection
//
// Config 包含状态命令所需的配置选项
// 提供执行状态检查时所需的依赖
type Config struct {
	Param       *newmigrate.MigrationParam // Migration connection // 迁移连接
	ScriptsPath string                     // Path to migration scripts DIR // 迁移脚本目录路径
	Objects     []any                      // GORM model objects used in schema comparison // 用于结构比较的 GORM 模型对象
}

// Status represents the current migration status
// Contains relevant information needed to understand migration state
//
// Status 表示当前迁移状态
// 包含理解迁移状态所需的相关信息
type Status struct {
	DatabaseVersion     uint     // Current database version // 当前数据库版本
	IsDirtyFlag         bool     // Database dirty state flag // 数据库脏状态标志
	HasMigrated         bool     // Migration applied flag // 迁移已应用标志
	LatestScriptVersion uint     // Latest version in scripts DIR // 脚本目录中的最新版本
	ScriptCount         int      // Count of migration scripts // 迁移脚本数量
	PendingCount        int      // Count of pending migrations // 待执行迁移数量
	HasUpToDate         bool     // Database is at latest version flag // 数据库已是最新版本标志
	PendingVersions     []uint   // List of pending migration versions // 待执行迁移版本列表
	SchemaDiffCount     int      // Count of schema differences // 结构差异数量
	SchemaDiffSQLs      []string // SQL statements showing schema differences // 结构差异的 SQL 语句
}

// GetStatus analyzes current migration state and returns comprehensive status
// Inspects database version, script versions and schema differences
// Returns Status struct containing relevant information
//
// GetStatus 分析当前迁移状态并返回综合状态
// 检查数据库版本、脚本版本和结构差异
// 返回包含相关信息的 Status 结构
func GetStatus(db *gorm.DB, migration *migrate.Migrate, scriptsPath string, objects []any) (*Status, error) {
	status := &Status{}

	// Get database version
	// 获取数据库版本
	version, dirtyFlag, err := migration.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			status.HasMigrated = false
			status.DatabaseVersion = 0
		} else {
			return nil, erero.Wro(err)
		}
	} else {
		status.HasMigrated = true
		status.DatabaseVersion = version
		status.IsDirtyFlag = dirtyFlag
	}

	// Scan scripts DIR and extract versions
	// 扫描脚本目录并提取版本
	scriptVersions, err := scanScriptVersions(scriptsPath)
	if err != nil {
		return nil, erero.Wro(err)
	}
	status.ScriptCount = len(scriptVersions)
	if len(scriptVersions) > 0 {
		status.LatestScriptVersion = scriptVersions[len(scriptVersions)-1]
	}

	// Calculate pending migrations
	// 计算待执行迁移
	for _, v := range scriptVersions {
		if v > status.DatabaseVersion {
			status.PendingVersions = append(status.PendingVersions, v)
		}
	}
	status.PendingCount = len(status.PendingVersions)
	status.HasUpToDate = status.PendingCount == 0

	// Check schema differences when objects are provided
	// 当提供对象时检查结构差异
	if len(objects) > 0 {
		diffSQLs := checkmigration.CheckMigrate(db, objects)
		status.SchemaDiffSQLs = diffSQLs
		status.SchemaDiffCount = len(diffSQLs)
	}

	return status, nil
}

// scanScriptVersions reads script DIR and extracts version numbers
// Returns sorted list of unique version numbers
//
// scanScriptVersions 读取脚本目录并提取版本号
// 返回排序后的唯一版本号列表
func scanScriptVersions(scriptsPath string) ([]uint, error) {
	entries, err := os.ReadDir(scriptsPath)
	if err != nil {
		return nil, erero.Wro(err)
	}

	versionSet := make(map[uint]bool)
	for _, item := range entries {
		if item.IsDir() {
			continue
		}
		migration, err := source.DefaultParse(item.Name())
		if err != nil {
			continue // Skip files that don't match migration pattern // 跳过不匹配迁移模式的文件
		}
		versionSet[migration.Version] = true
	}

	versions := make([]uint, 0, len(versionSet))
	for v := range versionSet {
		versions = append(versions, v)
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	return versions, nil
}

// ShowStatus outputs status information in a readable format
// Colored output improves reading experience
//
// ShowStatus 以可读格式输出状态信息
// 彩色输出提升阅读体验
func ShowStatus(status *Status) {
	eroticgo.CYAN.ShowMessage("=== Migration Status ===")

	// Database version
	// 数据库版本
	if status.HasMigrated {
		if status.IsDirtyFlag {
			eroticgo.RED.ShowMessage(fmt.Sprintf("Database Version: %d (DIRTY)", status.DatabaseVersion))
		} else {
			eroticgo.GREEN.ShowMessage(fmt.Sprintf("Database Version: %d", status.DatabaseVersion))
		}
	} else {
		eroticgo.YELLOW.ShowMessage("Database Version: (none - no migration records)")
	}

	// Script info
	// 脚本信息
	if status.ScriptCount > 0 {
		eroticgo.GREEN.ShowMessage(fmt.Sprintf("Scripts Latest: %d (%d scripts)", status.LatestScriptVersion, status.ScriptCount))
	} else {
		eroticgo.YELLOW.ShowMessage("Scripts Latest: (none - no scripts found)")
	}

	// Pending migrations
	// 待执行迁移
	if status.PendingCount > 0 {
		eroticgo.YELLOW.ShowMessage(fmt.Sprintf("Pending Migrations: %d", status.PendingCount))
		fmt.Println("  Versions:", status.PendingVersions)
	} else {
		eroticgo.GREEN.ShowMessage("Pending Migrations: 0 (up to date)")
	}

	// Schema differences
	// 结构差异
	if status.SchemaDiffCount > 0 {
		eroticgo.YELLOW.ShowMessage(fmt.Sprintf("Schema Differences: %d", status.SchemaDiffCount))
		fmt.Println("  (Database has changes not yet in migration scripts)")
		for i, sql := range status.SchemaDiffSQLs {
			fmt.Println("->", i+1, "->", sql)
		}
	} else if status.SchemaDiffCount == 0 && len(status.SchemaDiffSQLs) == 0 {
		eroticgo.GREEN.ShowMessage("Schema Differences: 0 (Models match database)")
	}
}

// NewStatusCmd creates cobra command that displays migration status
// Provides comprehensive view of current migration state
//
// NewStatusCmd 创建显示迁移状态的 cobra 命令
// 提供当前迁移状态的综合视图
func NewStatusCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		Long:  "Show current database version, script versions, pending migrations and schema differences",
		Run: func(cmd *cobra.Command, args []string) {
			migration, cleanup := cfg.Param.GetMigration()
			defer cleanup()

			db, cleanup2 := cfg.Param.GetDB()
			defer cleanup2()
			status := rese.P1(GetStatus(db, migration, cfg.ScriptsPath, cfg.Objects))
			ShowStatus(status)
		},
	}
}
