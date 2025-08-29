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

	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CustomLogger implements GORM logger interface to capture executed SQL statements
// Uses slice collection to record all SQL operations during DryRun mode
// Provides tracing capabilities for migration analysis and script generation
//
// CustomLogger 实现 GORM 日志接口来捕获执行的 SQL 语句
// 使用切片集合在 DryRun 模式下记录所有 SQL 操作
// 为迁移分析和脚本生成提供跟踪功能
type CustomLogger struct {
	SQLs []string // Collection of captured SQL statements // 捕获的 SQL 语句集合
}

func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	zaplog.SUG.Debugln("mode", int(level))
	return c
}

func (c *CustomLogger) Info(_ context.Context, msg string, data ...interface{}) {
	zaplog.SUG.Infoln("info", fmt.Sprintf(msg, data...))
}

func (c *CustomLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	zaplog.SUG.Warnln("warn", fmt.Sprintf(msg, data...))
}

func (c *CustomLogger) Error(_ context.Context, msg string, data ...interface{}) {
	zaplog.SUG.Errorln("error", fmt.Sprintf(msg, data...))
}

func (c *CustomLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	sqx, _ := fc()
	zaplog.SUG.Debugln("SQL>>>", eroticgo.GREEN.Sprint(sqx), "<<<END")
	c.SQLs = append(c.SQLs, sqx)
}

// GetMigrateOps analyzes GORM models and generates migration operations based on database differences
// Uses GORM DryRun mode with custom logger to capture SQL statements without execution
// Returns structured migration operations with both forward and reverse SQL scripts
//
// GetMigrateOps 分析 GORM 模型并基于数据库差异生成迁移操作
// 使用 GORM DryRun 模式和自定义日志来捕获 SQL 语句而不执行
// 返回包含正向和反向 SQL 脚本的结构化迁移操作
func GetMigrateOps(db *gorm.DB, objects []interface{}) MigrationOps {
	// Create custom logger for SQL capture
	// 创建用于 SQL 捕获的自定义日志
	customLogger := &CustomLogger{
		SQLs: make([]string, 0),
	}

	// Use DryRun mode with custom logger to capture SQL without execution
	// 使用 DryRun 模式和自定义日志来捕获 SQL 而不执行
	session := &gorm.Session{
		DryRun: true,
		Logger: customLogger,
	}
	must.Done(db.Session(session).AutoMigrate(objects...))

	// Display captured SQL statements for debugging
	// 显示捕获的 SQL 语句用于调试
	zaplog.SUG.Debugln("execute:", eroticgo.BLUE.Sprint(neatjsons.S(customLogger.SQLs)))

	results := make([]*MigrationOp, 0, len(customLogger.SQLs))
	for _, forwardSQL := range customLogger.SQLs {
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
