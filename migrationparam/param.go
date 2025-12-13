// Package migrationparam: Database migration instance factory with multiple initialization strategies
// Provides flexible migration creation supporting file systems, embedded resources and database drivers
//
// migrationparam: 数据库迁移实例工厂，支持多种初始化策略
// 提供灵活的迁移创建，支持文件系统、嵌入资源和数据库驱动
package migrationparam

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"gorm.io/gorm"
)

// MigrationParam provides unified database connection and migration management
// Contains factory functions for creating connections and migrations on demand
// Cleanup method handles proper resource release after operations complete
//
// MigrationParam 提供统一的数据库连接和迁移管理
// 包含按需创建连接和迁移的工厂函数
// Cleanup 方法在操作完成后处理资源释放
type MigrationParam struct {
	newDB        func() *gorm.DB // Factory to create database connection on demand // 按需创建数据库连接的工厂函数
	db           *gorm.DB
	newMigration func(db *gorm.DB) *migrate.Migrate // Factory that accepts shared database connection // 接受共享数据库连接的工厂函数
	migration    *migrate.Migrate
}

func NewMigrationParam(newDB func() *gorm.DB, newMigration func(db *gorm.DB) *migrate.Migrate) *MigrationParam {
	return &MigrationParam{
		newDB:        newDB,
		newMigration: newMigration,
	}
}

func (p *MigrationParam) GetDB() (*gorm.DB, func()) {
	if p.db == nil {
		p.db = p.newDB()
	}
	return p.db, p.cleanup
}

func (p *MigrationParam) GetMigration() (*migrate.Migrate, func()) {
	if p.migration == nil {
		db, cleanup := p.GetDB()
		p.migration = p.newMigration(db)
		_ = cleanup // 结果里的 cleanup 包含这个的 cleanup
	}
	return p.migration, p.cleanup
}

// cleanup releases resources after migration operations complete
// Closes migration instance and database connection properly
// Logs errors but continues cleanup to ensure resources are released
//
// cleanup 在迁移操作完成后释放资源
// 正确关闭迁移实例和数据库连接
// 记录错误但继续清理以确保资源被释放
func (p *MigrationParam) cleanup() {
	// Close migration instance
	// 关闭迁移实例
	if p.migration != nil {
		err1, err2 := p.migration.Close()
		must.Done(err1)
		must.Done(err2)
		p.migration = nil
	}

	// Close database connection
	// 关闭数据库连接
	if p.db != nil {
		must.Done(rese.P1(p.db.DB()).Close())
		p.db = nil
	}
}
