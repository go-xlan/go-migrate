// Package newmigrate: Database migration instance creation with multiple initialization strategies
// Provides flexible migration setup supporting file systems, embedded resources and database drivers
// Features generic type parameters and automatic registration and configuration
// Integrates with golang-migrate to provide robust versioning and execution
//
// newmigrate: 数据库迁移实例工厂，支持多种初始化策略
// 提供灵活的迁移创建，支持文件系统、嵌入资源和数据库驱动
// 具有泛型类型参数，用于自动驱动注册和配置
// 与 golang-migrate 库集成，提供稳健的版本控制和执行
package newmigrate

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func init() {
	must.Full(&file.File{}) // Register file source driver (side effects) // 注册文件源驱动（副作用）
}

// ScriptsAndDBSourceParam contains configuration using file-based migration with database connection string
// Supports migration from file system with database connection URL
// Used in simple setup scenarios when database connection string is available
//
// ScriptsAndDBSourceParam 包含基于文件的迁移配置和数据库连接字符串
// 支持从本地文件系统进行迁移，使用数据库连接 URL
// 用于可获得数据库连接字符串的简单设置场景
type ScriptsAndDBSourceParam struct {
	ScriptsInRoot string // Path to migration scripts DIR // 迁移脚本 DIR 路径
	ConnectSource string // Database connection string // 数据库连接字符串
}

// NewWithScriptsAndDBSource creates migration instance using file system scripts and database connection string
// Generic type param T enforces database driver interface compliance and triggers registration
// Supports multiple database types through golang-migrate driver system
// Returns configured migration instance prepared to execute
//
// Supported database drivers:
//   - sqlite3.Sqlite from "github.com/golang-migrate/migrate/v4/database/sqlite3"
//   - mysql.Mysql from "github.com/golang-migrate/migrate/v4/database/mysql"
//   - postgres.Postgres from "github.com/golang-migrate/migrate/v4/database/postgres"
//
// Usage examples:
//   - migration, err := NewWithScriptsAndDBSource[*sqlite3.Sqlite](param)
//   - migration, err := NewWithScriptsAndDBSource[*mysql.Mysql](param)
//   - migration, err := NewWithScriptsAndDBSource[*postgres.Postgres](param)
//
// NewWithScriptsAndDBSource 使用文件系统脚本和数据库连接字符串创建迁移实例
// 泛型类型参数 T 强制数据库驱动接口兼容性并触发驱动注册
// 通过 golang-migrate 驱动系统支持多种数据库类型
// 返回已配置的迁移实例，可用于执行
func NewWithScriptsAndDBSource[T database.Driver](param *ScriptsAndDBSourceParam) (*migrate.Migrate, error) {
	sourceURL := "file://" + param.ScriptsInRoot
	migration, err := migrate.New(
		sourceURL,
		param.ConnectSource,
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}

// ScriptsAndDatabaseParam contains configuration for file-based migration with database driver instance
// Provides direct database driver management in advanced configuration scenarios
// Enables custom database setup and connection management
//
// ScriptsAndDatabaseParam 包含基于文件的迁移配置和数据库驱动实例
// 为高级配置场景提供直接的数据库驱动控制
// 支持自定义数据库设置和连接管理
type ScriptsAndDatabaseParam struct {
	ScriptsInRoot    string          // Path to migration scripts DIR // 迁移脚本 DIR 路径
	DatabaseName     string          // Database name ID // 数据库名称标识
	DatabaseInstance database.Driver // Database driver instance // 数据库驱动实例
}

// NewWithScriptsAndDatabase creates migration instance using file system scripts and database driver instance
// Provides direct database driver management with file-based migration scripts
// Returns configured migration instance prepared to execute
//
// NewWithScriptsAndDatabase 使用文件系统脚本和数据库驱动实例创建迁移实例
// 通过基于文件的迁移脚本提供直接的数据库驱动控制
// 返回已配置的迁移实例，可运行
func NewWithScriptsAndDatabase(param *ScriptsAndDatabaseParam) (*migrate.Migrate, error) {
	sourceURL := "file://" + param.ScriptsInRoot
	migration, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		param.DatabaseName,
		param.DatabaseInstance,
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}

// EmbedFsAndDatabaseParam contains configuration for embedded file system migration with database driver
// Enables migration scripts to be embedded into executable file
// Supports self-contained applications with built-in migration capabilities
//
// EmbedFsAndDatabaseParam 包含嵌入文件系统迁移的配置和数据库驱动
// 支持将迁移脚本嵌入到二进制文件中进行分发
// 支持具有内置迁移功能的自包含应用程序
type EmbedFsAndDatabaseParam struct {
	MigrationsFS     *embed.FS       // Embedded file system with migrations // 包含迁移的嵌入文件系统
	EmbedDirName     string          // DIR name within embedded FS // 嵌入 FS 中的 DIR 名称
	DatabaseName     string          // Database name ID // 数据库名称标识
	DatabaseInstance database.Driver // Database driver instance // 数据库驱动实例
}

// NewWithEmbedFsAndDatabase creates migration instance using embedded file system and database driver
// Enables self-contained binaries with built-in migration scripts
// Returns configured migration instance prepared to execute
//
// NewWithEmbedFsAndDatabase 使用嵌入文件系统和数据库驱动创建迁移实例
// 支持带有内置迁移脚本的自包含二进制文件
// 返回已配置的迁移实例，可运行
func NewWithEmbedFsAndDatabase(param *EmbedFsAndDatabaseParam) (*migrate.Migrate, error) {
	const sourceName = "iofs"
	// Reference: https://github.com/golang-migrate/migrate/blob/278833935c12dda022b1355f33a897d895501c45/source/iofs/example_test.go#L22
	// 详情参考: https://github.com/golang-migrate/migrate/blob/278833935c12dda022b1355f33a897d895501c45/source/iofs/example_test.go#L22
	migration, err := migrate.NewWithInstance(
		sourceName, // Fixed iofs type // 固定的 iofs 类型
		rese.V1(iofs.New(param.MigrationsFS, param.EmbedDirName)), // Initialize iofs driver // 初始化 iofs 驱动
		param.DatabaseName,
		param.DatabaseInstance,
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return migration, nil
}
