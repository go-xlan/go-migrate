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

// CustomLogger 自定义 Logger 捕获执行过的 SQL 语句
type CustomLogger struct {
	SQLs []string
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

func GetMigrateOps(db *gorm.DB, objects []interface{}) MigrationOps {
	// 创建自定义 Logger
	customLogger := &CustomLogger{
		SQLs: make([]string, 0),
	}

	// 使用 DryRun 和自定义 Logger
	session := &gorm.Session{
		DryRun: true,
		Logger: customLogger,
	}
	must.Done(db.Session(session).AutoMigrate(objects...))

	// 获取生成的 SQL
	zaplog.SUG.Debugln("execute:", eroticgo.BLUE.Sprint(neatjsons.S(customLogger.SQLs)))

	results := make([]*MigrationOp, 0, len(customLogger.SQLs))
	for _, forwardSQL := range customLogger.SQLs {
		// 判断是否需要迁移
		if migrationOp, match := NewMigrationOp(forwardSQL); match {
			results = append(results, must.Full(migrationOp))
		}
	}
	return results
}

// CheckMigrate 这个函数用来对比DB和模型，检查DB里缺少什么没有执行的语句
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
