package checkmigration

import (
	"fmt"
	"slices"
	"strings"

	"github.com/yyle88/tern"
	"github.com/yyle88/tern/zerotern"
)

const (
	// raiseStatement provides safe placeholder for unimplemented reverse migrations
	// Prevents accidental execution while clearly marking incomplete implementations
	// Note: "SELECT 1 / 0;" was tested but doesn't effectively cause script failure
	//
	// raiseStatement 为未实现的反向迁移提供安全占位符
	// 防止意外执行，同时明确标记不完整实现
	// 注意：测试过 "SELECT 1 / 0;" 但不能有效地导致脚本失败
	raiseStatement = "SELECT TODO / PANIC / RAISE / THROW;"
)

// MigrationKind represents the type and characteristics of a database migration operation
// Contains substring patterns for identifying operation type and corresponding reverse operation
// Used for automated reverse script generation and operation categorization
//
// MigrationKind 表示数据库迁移操作的类型和特征
// 包含用于识别操作类型和对应反向操作的子字符串模式
// 用于自动反向脚本生成和操作分类
type MigrationKind struct {
	ForwardSubstr string // SQL pattern for forward operation // 正向操作的 SQL 模式
	ReverseSubstr string // SQL pattern for reverse operation // 反向操作的 SQL 模式
}

var migrationKinds = []*MigrationKind{
	{ForwardSubstr: "CREATE TABLE", ReverseSubstr: "DROP TABLE"},
	{ForwardSubstr: "ALTER TABLE", ReverseSubstr: "ALTER TABLE"},
	{ForwardSubstr: "ADD COLUMN", ReverseSubstr: "DROP COLUMN"},
	{ForwardSubstr: "ADD INDEX", ReverseSubstr: "DROP INDEX"},
	{ForwardSubstr: "CREATE INDEX", ReverseSubstr: "DROP INDEX"},
	{ForwardSubstr: "CREATE UNIQUE INDEX", ReverseSubstr: "DROP INDEX"},
}

// MigrationOp represents a single database migration operation with forward SQL and operation metadata
// Combines SQL statement with operation type information for comprehensive migration management
// Supports both forward execution and reverse script generation
//
// MigrationOp 表示单个数据库迁移操作，包含正向 SQL 和操作元数据
// 将 SQL 语句与操作类型信息相结合，用于全面的迁移管理
// 支持正向执行和反向脚本生成
type MigrationOp struct {
	ForwardSQL string         // SQL statement for forward migration // 正向迁移的 SQL 语句
	Kind       *MigrationKind // Operation type and reverse pattern // 操作类型和反向模式
}

// NewMigrationOp creates migration operation from SQL statement by matching against known patterns
// Analyzes SQL content to determine operation type and appropriate reverse operation
// Returns migration operation instance and success flag
//
// NewMigrationOp 通过与已知模式匹配，从 SQL 语句创建迁移操作
// 分析 SQL 内容来确定操作类型和适当的反向操作
// 返回迁移操作实例和成功标志
func NewMigrationOp(forwardSQL string) (*MigrationOp, bool) {
	for _, sub := range migrationKinds {
		if strings.Contains(forwardSQL, sub.ForwardSubstr) {
			return &MigrationOp{
				ForwardSQL: forwardSQL,
				Kind: &MigrationKind{ // clone it and return to outside
					ForwardSubstr: sub.ForwardSubstr,
					ReverseSubstr: sub.ReverseSubstr,
				},
			}, true
		}
	}
	return nil, false
}

func (op *MigrationOp) GetForwardSQL() string {
	return op.ForwardSQL
}

func (op *MigrationOp) GetReverseSQL() (string, bool) {
	return raiseStatement + " -- " + op.Kind.ReverseSubstr, false //认为肯定有更好的工具来做这件事，这里暂时不要实现这个功能(别整太复杂反正早期也没人用的)
}

// MigrationOps represents a collection of migration operations with batch processing capabilities
// Provides methods for searching, SQL extraction, and script generation
// Supports both forward and reverse migration script creation
//
// MigrationOps 表示具有批量处理功能的迁移操作集合
// 提供搜索、SQL 提取和脚本生成方法
// 支持正向和反向迁移脚本创建
type MigrationOps []*MigrationOp

func (ops MigrationOps) SearchOp(forwardSQL string) *MigrationOp {
	posIndex := slices.IndexFunc(ops, func(op *MigrationOp) bool {
		return strings.EqualFold(op.ForwardSQL, forwardSQL)
	})
	return tern.BF(posIndex >= 0, func() *MigrationOp {
		return ops[posIndex]
	})
}

func (ops MigrationOps) GetForwardSQLs() []string {
	var sqs = make([]string, 0, len(ops))
	for _, op := range ops {
		sqs = append(sqs, op.ForwardSQL)
	}
	return sqs
}

func (ops MigrationOps) GetForwardScript() string {
	var sqs = make([]string, 0, len(ops))
	for _, op := range ops {
		sqs = append(sqs, op.GetForwardSQL()+";")
	}
	res := strings.Join(sqs, "\n\n")
	if len(res) > 0 {
		res += "\n"
	}
	return res
}

func (ops MigrationOps) GetReverseScript() (string, bool) {
	var sqs = make([]string, 0, len(ops))
	var okk = true
	// 需要倒序执行逆向的操作，比如先删除索引再删除列，这样拼接出来
	for idx := len(ops) - 1; idx >= 0; idx-- {
		op := ops[idx]
		reverseSQL, ok := op.GetReverseSQL()
		if ok {
			sqs = append(sqs, reverseSQL+";")
			continue
		}
		okk = false //只要有一个出错，就记录下来有错
		reverseSQL = zerotern.VF(reverseSQL, func() string {
			return raiseStatement + " -- " + op.Kind.ReverseSubstr
		})
		forwardSQL := op.GetForwardSQL()
		sqLine := fmt.Sprintf("-- reverse -- %s;\n%s; -- TODO", forwardSQL, reverseSQL)
		sqs = append(sqs, sqLine) //当有错误时，就使用注释“待做”拼接它们
	}
	res := strings.Join(sqs, "\n\n")
	if len(res) > 0 {
		res += "\n"
	}
	return res, okk
}
