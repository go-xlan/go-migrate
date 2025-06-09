package checkmigration_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractTableNameFromCreateTable(t *testing.T) {
	t.Run("case-1", func(t *testing.T) {
		sql := "CREATE TABLE `users` (`id` integer PRIMARY KEY AUTOINCREMENT,`name` text,`code` text,CONSTRAINT `uni_users_code` UNIQUE (`code`))"
		tableName := extractTableNameFromCreateTable(sql)
		t.Log(tableName)
		require.Equal(t, "users", tableName)
	})
	t.Run("case-2", func(t *testing.T) {
		sql := "CREATE TABLE `infos` (`code` varchar(255) COMMENT '编码(通用缩写)', `name` varchar(255) COMMENT '名称(英文名称)', PRIMARY KEY (`code`))"
		tableName := extractTableNameFromCreateTable(sql)
		t.Log(tableName)
		require.Equal(t, "infos", tableName)
	})
}

func Test_extractTableAndColumnFromAlterTable(t *testing.T) {
	t.Run("case-1", func(t *testing.T) {
		table, col := extractTableAndColumnFromAlterTableAddColune("ALTER TABLE `users` ADD `age` bigint")
		require.Equal(t, "users", table)
		require.Equal(t, "age", col)
	})

	t.Run("case-2", func(t *testing.T) {
		table, col := extractTableAndColumnFromAlterTableAddColune("ALTER TABLE `users` ADD `from` varchar(255)")
		require.Equal(t, "users", table)
		require.Equal(t, "from", col)
	})

	t.Run("case-3", func(t *testing.T) {
		table, col := extractTableAndColumnFromAlterTableAddColune("ALTER TABLE `users` ADD `student_no` varchar(255)")
		require.Equal(t, "users", table)
		require.Equal(t, "student_no", col)
	})

	t.Run("case-4", func(t *testing.T) {
		table, col := extractTableAndColumnFromAlterTableAddColune("ALTER TABLE `users` ADD `rank` integer")
		require.Equal(t, "users", table)
		require.Equal(t, "rank", col)
	})
}

func Test_extractIndexAndTableFromCreateIndex(t *testing.T) {
	t.Run("case-1", func(t *testing.T) {
		index, table := extractIndexAndTableFromCreateIndex("CREATE UNIQUE INDEX `idx_users_student_no` ON `users`(`student_no`)")
		require.Equal(t, "idx_users_student_no", index)
		require.Equal(t, "users", table)
	})

	t.Run("case-2", func(t *testing.T) {
		index, table := extractIndexAndTableFromCreateIndex("CREATE INDEX `idx_users_rank` ON `users`(`rank`)")
		require.Equal(t, "idx_users_rank", index)
		require.Equal(t, "users", table)
	})
}

func extractTableNameFromCreateTable(sql string) string {
	re := regexp.MustCompile(`(?i)^CREATE TABLE\s+` + "`" + `(\w+)` + "`" + `\s*\(.*\)$`)
	match := re.FindStringSubmatch(sql)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

func extractTableAndColumnFromAlterTableAddColune(sql string) (table string, column string) {
	re := regexp.MustCompile(`(?i)^ALTER TABLE\s+` + "`" + `(\w+)` + "`" + `\s+ADD(?: COLUMN)?\s+` + "`" + `(\w+)` + "`" + `\s+\w+`)
	match := re.FindStringSubmatch(sql)
	if len(match) >= 3 {
		return match[1], match[2]
	}
	return "", ""
}

func extractIndexAndTableFromCreateIndex(sql string) (index string, table string) {
	re := regexp.MustCompile(`(?i)^CREATE(?: UNIQUE)? INDEX\s+` + "`" + `(\w+)` + "`" + `\s+ON\s+` + "`" + `(\w+)` + "`")
	match := re.FindStringSubmatch(sql)
	if len(match) >= 3 {
		return match[1], match[2]
	}
	return "", ""
}
