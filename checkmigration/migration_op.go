package checkmigration

import (
	"fmt"
	"slices"
	"strings"

	"github.com/yyle88/tern"
	"github.com/yyle88/tern/zerotern"
)

type MigrationKind struct {
	ForwardSubstr string
	ReverseSubstr string
}

var migrationKinds = []*MigrationKind{
	{ForwardSubstr: "CREATE TABLE", ReverseSubstr: "DROP TABLE"},
	{ForwardSubstr: "ALTER TABLE", ReverseSubstr: "ALTER TABLE"},
	{ForwardSubstr: "ADD COLUMN", ReverseSubstr: "DROP COLUMN"},
	{ForwardSubstr: "ADD INDEX", ReverseSubstr: "DROP INDEX"},
	{ForwardSubstr: "CREATE INDEX", ReverseSubstr: "DROP INDEX"},
	{ForwardSubstr: "CREATE UNIQUE INDEX", ReverseSubstr: "DROP INDEX"},
}

type MigrationOp struct {
	ForwardSQL string
	Kind       *MigrationKind
}

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
	return "SELECT 1 / 0; -- " + op.Kind.ReverseSubstr, false //认为肯定有更好的工具来做这件事，这里暂时不要实现这个功能(别整太复杂反正早期也没人用的)
}

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
			return "SELECT 1 / 0; -- " + op.Kind.ReverseSubstr
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
