// Package utils: Internal utility functions for migration operations and error handling
// Provides common helper functions for UUID generation and migration result processing
// Includes specialized error handling for golang-migrate specific error cases
//
// utils: 用于迁移操作和错误处理的内部工具函数
// 提供 UUID 生成和迁移结果处理的通用助手函数
// 包含针对 golang-migrate 特定错误情况的专门错误处理
package utils

import (
	"encoding/hex"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/zaplog"
)

// NewUUID32s generates 32-character hexadecimal UUID string for unique identification
// Creates standard UUID and converts to lowercase hex representation
// Used for generating unique identifiers in migration contexts
//
// NewUUID32s 生成 32 字符十六进制 UUID 字符串用于唯一标识
// 创建标准 UUID 并转换为小写十六进制表示
// 用于在迁移上下文中生成唯一标识符
func NewUUID32s() string {
	u := uuid.New()
	return hex.EncodeToString(u[:])
}

// WhistleCause processes migration errors with appropriate logging and panic behavior
// Handles common golang-migrate error cases with informative messages
// Uses color-coded output for different error types and success states
//
// WhistleCause 处理迁移错误，采用适当的日志和异常行为
// 处理常见的 golang-migrate 错误情况，并提供信息性消息
// 使用颜色编码输出来区分不同的错误类型和成功状态
func WhistleCause(cause error) {
	if cause != nil {
		if errors.Is(cause, migrate.ErrNoChange) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("NO MIGRATION FILES TO RUN"))
		} else if errors.Is(cause, migrate.ErrNilVersion) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("NO VERSION IN VERSION-TABLE(schema_migrations)"))
		} else if errors.Is(cause, os.ErrNotExist) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("MIGRATION FILES NOT FOUND"))
		} else {
			zaplog.SUG.Panicln(eroticgo.RED.Sprint("MIGRATION FAILED:"), cause)
		}
		return
	}
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint("MIGRATION SUCCESS"))
}
