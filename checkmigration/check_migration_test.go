package checkmigration_test

import (
	"strings"
	"testing"

	"github.com/go-xlan/go-migrate/checkmigration"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

func TestMain(m *testing.M) {
	db := done.VCE(gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})).Nice()
	defer func() {
		must.Done(rese.P1(db.DB()).Close())
	}()

	caseDB = db
	m.Run()
}

func TestCheckMigrate(t *testing.T) {
	db := caseDB

	{
		migrateSQLs := checkmigration.CheckMigrate(db, []any{&UserV1{}})
		require.Len(t, migrateSQLs, 1)
		createTable := migrateSQLs[0]
		// 这里只需要确认是1个create table语句
		require.True(t, strings.HasPrefix(createTable, "CREATE TABLE `users` ("))
		require.True(t, strings.HasSuffix(createTable, ")"))

		require.NoError(t, db.AutoMigrate(&UserV1{}))
	}

	{
		migrateSQLs := checkmigration.CheckMigrate(db, []any{&UserV2{}})
		require.Len(t, migrateSQLs, 6)
		//因为这里并不检查顺序所以使用 contains 断言
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `age` bigint")
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `from` varchar(255)")
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `student_no` varchar(255)")
		require.Contains(t, migrateSQLs, "ALTER TABLE `users` ADD `rank` integer")
		require.Contains(t, migrateSQLs, "CREATE UNIQUE INDEX `idx_users_student_no` ON `users`(`student_no`)")
		require.Contains(t, migrateSQLs, "CREATE INDEX `idx_users_rank` ON `users`(`rank`)")

		require.NoError(t, db.AutoMigrate(&UserV2{}))
	}
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
