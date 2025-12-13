// Package checkmigration: Database schema migration validation and SQL generation engine
// Performs intelligent comparison between GORM model definitions and actual database schemas
// Captures and parses SQL statements from GORM's DryRun mode to identify migration requirements
// Generates both forward and reverse migration scripts with proper categorization
//
// checkmigration: 数据库结构迁移验证和 SQL 生成引擎
// 在 GORM 模型定义与实际数据库结构之间执行智能比较
// 从 GORM 的 DryRun 模式捕获和解析 SQL 语句，识别迁移需求
// 生成正向和反向迁移脚本，并进行适当分类
package checkmigration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-xlan/go-migrate/migrationparam"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SqlCapture implements GORM logger interface to capture executed SQL statements
// Uses slice collection to record SQL operations during DryRun mode
// Provides tracing capabilities for migration analysis and script generation
//
// SqlCapture 实现 GORM 日志接口来捕获执行的 SQL 语句
// 使用切片集合在 DryRun 模式下记录 SQL 操作
// 为迁移分析和脚本生成提供跟踪功能
type SqlCapture struct {
	SQLs      []string // Collection of captured SQL statements // 捕获的 SQL 语句集合
	debugMode bool     // Enable debug output // 调试模式
}

func (c *SqlCapture) LogMode(level logger.LogLevel) logger.Interface {
	if c.debugMode {
		zaplog.SUG.Debugln("mode", int(level))
	}
	return c
}

func (c *SqlCapture) Info(_ context.Context, msg string, data ...interface{}) {
	if c.debugMode {
		zaplog.SUG.Infoln("info", fmt.Sprintf(msg, data...))
	}
}

func (c *SqlCapture) Warn(_ context.Context, msg string, data ...interface{}) {
	if c.debugMode {
		zaplog.SUG.Warnln("warn", fmt.Sprintf(msg, data...))
	}
}

func (c *SqlCapture) Error(_ context.Context, msg string, data ...interface{}) {
	if c.debugMode {
		zaplog.SUG.Errorln("error", fmt.Sprintf(msg, data...))
	}
}

func (c *SqlCapture) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	sqx, _ := fc()
	if c.debugMode {
		zaplog.SUG.Debugln("SQL>>>", eroticgo.GREEN.Sprint(sqx), "<<<END")
	}
	c.SQLs = append(c.SQLs, sqx)
}

// GetMigrateOps analyzes GORM models and generates migration operations based on database differences
// Uses GORM DryRun mode with custom logger to capture SQL statements without execution
// Returns structured migration operations with both forward and reverse SQL scripts
// Debug output controlled by package-level SetDebugMode
//
// GetMigrateOps 分析 GORM 模型并基于数据库差异生成迁移操作
// 使用 GORM DryRun 模式和自定义日志来捕获 SQL 语句而不执行
// 返回包含正向和反向 SQL 脚本的结构化迁移操作
// 调试输出由包级别的 SetDebugMode 控制
func GetMigrateOps(db *gorm.DB, objects []interface{}) MigrationOps {
	// Create SqlCapture for SQL capture
	// 创建 SqlCapture 用于 SQL 捕获
	sqlCapture := &SqlCapture{
		SQLs:      make([]string, 0),
		debugMode: migrationparam.GetDebugMode(),
	}

	// Use DryRun mode with SqlCapture to capture SQL without execution
	// 使用 DryRun 模式和 SqlCapture 来捕获 SQL 而不执行
	session := &gorm.Session{
		DryRun: true,
		Logger: sqlCapture,
	}
	must.Done(db.Session(session).AutoMigrate(objects...))

	// Display captured SQL statements for debugging
	// 显示捕获的 SQL 语句用于调试
	if sqlCapture.debugMode {
		zaplog.SUG.Debugln("execute:", eroticgo.BLUE.Sprint(neatjsons.S(sqlCapture.SQLs)))
	}

	results := make([]*MigrationOp, 0, len(sqlCapture.SQLs))
	for _, forwardSQL := range sqlCapture.SQLs {
		// Parse SQL to determine if migration is needed
		// 解析 SQL 以确定是否需要迁移
		if migrationOp, match := NewMigrationOp(forwardSQL); match {
			results = append(results, must.Full(migrationOp))
		}
	}
	return results
}

// CheckMigrate compares database schema against GORM models and returns missing SQL statements
// Performs comprehensive analysis to identify required database migrations
// Returns list of forward SQL statements that need to be executed
//
// CheckMigrate 对比数据库结构与 GORM 模型，返回缺失的 SQL 语句
// 执行全面分析来识别所需的数据库迁移
// 返回需要执行的正向 SQL 语句列表
func CheckMigrate(db *gorm.DB, objects []interface{}) []string {
	steps := GetMigrateOps(db, objects)
	zaplog.LOG.Debug("missing", zap.Int("size", len(steps)))
	sqs := steps.GetForwardSQLs()
	if len(sqs) > 0 {
		debugMigrationSqs(sqs)
	}
	zaplog.SUG.Debugln("success")
	return sqs
}

func debugMigrationSqs(sqs []string) {
	zaplog.SUG.Debugln("-")
	for idx, sqx := range sqs {
		zaplog.SUG.Debug(
			"missing:",
			fmt.Sprintf("(%d/%d)", idx, len(sqs)),
			"\n",
			eroticgo.PINK.Sprint("----------------"),
			"\n\n",
			eroticgo.PINK.Sprint(sqx+";"),
			"\n\n",
			eroticgo.PINK.Sprint("----------------"),
		)
	}
	zaplog.SUG.Debugln("-")
	zaplog.SUG.Debugln("-")
	zaplog.SUG.Debug(
		"scripts:",
		"\n",
		eroticgo.CYAN.Sprint("----------------"),
		"\n\n",
		eroticgo.CYAN.Sprint(strings.Join(sqs, ";\n\n")+";"),
		"\n\n",
		eroticgo.CYAN.Sprint("----------------"),
	)
	zaplog.SUG.Debugln("-")
	zaplog.SUG.Debugln("-")
}
