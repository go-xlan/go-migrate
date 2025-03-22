package checkmigration_test

import (
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
	db := done.VCE(gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
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

	require.NotEmpty(t, checkmigration.CheckMigrate(db, []any{&UserV1{}}))

	require.NoError(t, db.AutoMigrate(&UserV1{}))

	require.NotEmpty(t, checkmigration.CheckMigrate(db, []any{&UserV2{}}))
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
