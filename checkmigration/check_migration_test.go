// Package checkmigration_test provides comprehensive tests that validate migration detection and SQL generation
// Uses in-memory SQLite database to test schema comparison without external dependencies
//
// checkmigration_test 包提供验证迁移检测和 SQL 生成的综合测试
// 使用内存 SQLite 数据库来测试结构比较，无需外部依赖
package checkmigration_test

import (
	"fmt"
	"testing"

	"github.com/go-xlan/go-migrate/checkmigration"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

// TestMain initializes shared in-memory database instance and runs all tests
//
// TestMain 初始化共享的内存数据库实例并运行所有测试
func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := done.VCE(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})).Nice()
	defer func() {
		must.Done(rese.P1(db.DB()).Close())
	}()

	caseDB = db
	m.Run()
}

// TestCheckMigrate validates migration SQL detection with incremental schema changes
// Tests sequential model upgrades from V1 to V2 with column and index additions
//
// TestCheckMigrate 验证带有增量结构变更的迁移 SQL 检测
// 测试从 V1 到 V2 的顺序模型升级，包括列和索引的添加
func TestCheckMigrate(t *testing.T) {
	db := caseDB

	require.True(t, t.Run("case-1", func(t *testing.T) {
		migrateSQLs := checkmigration.CheckMigrate(db, []any{&UserV1{}})
		require.Len(t, migrateSQLs, 1)
		// Confirm single CREATE TABLE statement // 确认是1个 CREATE TABLE 语句
		tableName := extractTableNameFromCreateTable(migrateSQLs[0])
		require.Equal(t, "users", tableName)

		require.NoError(t, db.AutoMigrate(&UserV1{}))
	}))

	require.True(t, t.Run("case-2", func(t *testing.T) {
		migrateSQLs := checkmigration.CheckMigrate(db, []any{&UserV2{}})
		require.Len(t, migrateSQLs, 6)
		// Use contains assertion since sequence is not guaranteed // 因为不检查顺序所以使用 contains 断言
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `age` bigint")
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `from` varchar(255)")
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `student_no` varchar(255)")
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `rank` integer")
		require.Contains(t, migrateSQLs, "CREATE UNIQUE INDEX `idx_users_student_no` ON `users`(`student_no`)")
		require.Contains(t, migrateSQLs, "CREATE INDEX `idx_users_rank` ON `users`(`rank`)")

		require.NoError(t, db.AutoMigrate(&UserV2{}))
	}))
}

type UserV1 struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
	Code string `gorm:"unique;"`
}

func (u *UserV1) TableName() string {
	return "users"
}

type UserV2 struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:200"`
	Age       int    `gorm:"type:bigint"`
	From      string `gorm:"type:varchar(255)"`
	StudentNo string `gorm:"type:varchar(255);uniqueIndex;"`
	Rank      int    `gorm:"type:int;index;"`
}

func (u *UserV2) TableName() string {
	return "users"
}

// TestCheckMigrate_Product validates migration detection with product models
// Tests CREATE TABLE, ALTER TABLE, and CREATE INDEX operations
//
// TestCheckMigrate_Product 验证产品模型的迁移检测
// 测试 CREATE TABLE、ALTER TABLE 和 CREATE INDEX 操作
func TestCheckMigrate_Product(t *testing.T) {
	db := caseDB

	require.True(t, t.Run("case-1", func(t *testing.T) {
		migrateOps := checkmigration.GetMigrateOps(db, []any{&ProductV1{}})
		require.Len(t, migrateOps, 1)
		op := migrateOps[0]
		tableName := extractTableNameFromCreateTable(op.ForwardSQL)
		require.Equal(t, "products", tableName)

		require.Equal(t, "CREATE TABLE", op.Kind.ForwardSubstr)
		require.Equal(t, "DROP TABLE", op.Kind.ReverseSubstr)

		showDebugScripts(t, migrateOps)

		require.NoError(t, db.AutoMigrate(&ProductV1{}))
	}))

	require.True(t, t.Run("case-2", func(t *testing.T) {
		migrateOps := checkmigration.GetMigrateOps(db, []any{&ProductV2{}})
		require.Len(t, migrateOps, 3)
		{
			op := requireOperation(t, migrateOps, "ALTER TABLE `products` ADD `price` float64")
			require.Equal(t, "ALTER TABLE", op.Kind.ForwardSubstr)
			require.Equal(t, "ALTER TABLE", op.Kind.ReverseSubstr)

			table, column := extractTableAndColumnFromAlterTableAddColune(op.ForwardSQL)
			require.Equal(t, "products", table)
			require.Equal(t, "price", column)
		}
		{
			op := requireOperation(t, migrateOps, "ALTER TABLE `products` ADD `sku` varchar(50)")
			require.Equal(t, "ALTER TABLE", op.Kind.ForwardSubstr)
			require.Equal(t, "ALTER TABLE", op.Kind.ReverseSubstr)

			table, column := extractTableAndColumnFromAlterTableAddColune(op.ForwardSQL)
			require.Equal(t, "products", table)
			require.Equal(t, "sku", column)
		}
		{
			op := requireOperation(t, migrateOps, "CREATE UNIQUE INDEX `idx_products_sku` ON `products`(`sku`)")
			require.Equal(t, "CREATE UNIQUE INDEX", op.Kind.ForwardSubstr)
			require.Equal(t, "DROP INDEX", op.Kind.ReverseSubstr)

			indexName, table := extractIndexAndTableFromCreateIndex(op.ForwardSQL)
			require.Equal(t, "products", table)
			require.Equal(t, "idx_products_sku", indexName)
		}

		showDebugScripts(t, migrateOps)

		require.NoError(t, db.AutoMigrate(&ProductV2{}))
	}))

	require.True(t, t.Run("case-3", func(t *testing.T) {
		migrateOps := checkmigration.GetMigrateOps(db, []any{&ProductV3{}})
		require.Len(t, migrateOps, 6)

		must.Full(requireOperation(t, migrateOps, "ALTER TABLE `products` ADD `brand` varchar(100)"))
		must.Full(requireOperation(t, migrateOps, "ALTER TABLE `products` ADD `country` varchar(100)"))
		{
			op := requireOperation(t, migrateOps, "CREATE INDEX `idx_brand_country_union` ON `products`(`brand`,`country`)")
			indexName, table := extractIndexAndTableFromCreateIndex(op.ForwardSQL)
			require.Equal(t, "products", table)
			require.Equal(t, "idx_brand_country_union", indexName)
		}

		must.Full(requireOperation(t, migrateOps, "ALTER TABLE `products` ADD `supplier_code` varchar(100)"))
		must.Full(requireOperation(t, migrateOps, "ALTER TABLE `products` ADD `batch_no` varchar(100)"))
		{
			op := requireOperation(t, migrateOps, "CREATE UNIQUE INDEX `ux_supplier_batch` ON `products`(`supplier_code`,`batch_no`)")
			indexName, table := extractIndexAndTableFromCreateIndex(op.ForwardSQL)
			require.Equal(t, "products", table)
			require.Equal(t, "ux_supplier_batch", indexName)
		}

		showDebugScripts(t, migrateOps)

		require.NoError(t, db.AutoMigrate(&ProductV3{}))
	}))
}

// showDebugScripts outputs forward and reverse scripts with colored formatting
//
// showDebugScripts 输出带颜色格式化的正向和反向脚本
func showDebugScripts(t *testing.T, migrateOps checkmigration.MigrationOps) {
	forwardScript := migrateOps.GetForwardScript()
	zaplog.ZAPS.Skip(1).SUG.Debug("forward:", "\n", eroticgo.AQUA.Sprint(forwardScript))
	reverseScript, _ := migrateOps.GetReverseScript()
	zaplog.ZAPS.Skip(1).SUG.Debug("reverse:", "\n", eroticgo.PINK.Sprint(reverseScript))
}

// requireOperation asserts operation exists and returns it with debug output
//
// requireOperation 断言操作存在并返回它，同时输出调试信息
func requireOperation(t *testing.T, migrateOps checkmigration.MigrationOps, forwardSQL string) *checkmigration.MigrationOp {
	t.Log(forwardSQL)
	op := migrateOps.SearchOp(forwardSQL)
	require.NotNil(t, op)
	t.Log(neatjsons.S(op))
	return op
}

type ProductV1 struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:200"`
}

func (p *ProductV1) TableName() string {
	return "products"
}

type ProductV2 struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"size:255"`
	Price float64 `gorm:"type:float64"`
	SKU   string  `gorm:"type:varchar(50);uniqueIndex"`
}

func (p *ProductV2) TableName() string {
	return "products"
}

type ProductV3 struct {
	ID           uint    `gorm:"primaryKey"`
	Name         string  `gorm:"size:255"`
	Price        float64 `gorm:"type:float64"`
	SKU          string  `gorm:"type:varchar(50);uniqueIndex"`
	Brand        string  `gorm:"type:varchar(100);index:idx_brand_country_union"` // 普通复合索引
	Country      string  `gorm:"type:varchar(100);index:idx_brand_country_union"` // 普通复合索引
	SupplierCode string  `gorm:"type:varchar(100);uniqueIndex:ux_supplier_batch"` // 唯一复合索引
	BatchNo      string  `gorm:"type:varchar(100);uniqueIndex:ux_supplier_batch"` // 唯一复合索引
}

func (p *ProductV3) TableName() string {
	return "products"
}
