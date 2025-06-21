package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/go-xlan/go-migrate/cobramigration"
	"github.com/go-xlan/go-migrate/internal/demos/demo1x/internal/models"
	"github.com/go-xlan/go-migrate/newmigrate"
	"github.com/go-xlan/go-migrate/newscripts"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// go run main.go next-script create --version-type TIME --description create_table
// go run main.go migrate all
func main() {
	db := newGormDB(&MysqlConfig{
		Source:   "tcp(localhost:3306)/xlan_migrate_demo1x?charset=utf8mb4&parseTime=true",
		Username: "root",
		Password: "123456",
	})
	defer rese.F0(rese.P1(db.DB()).Close)

	scriptsInRoot := runpath.PARENT.Join("scripts")

	migration := rese.P1(newmigrate.NewWithScriptsAndDatabase(
		&newmigrate.ScriptsAndDatabaseParam{
			ScriptsInRoot:    scriptsInRoot,
			DatabaseName:     "mysql",
			DatabaseInstance: rese.V1(mysqlmigrate.WithInstance(rese.P1(db.DB()), &mysqlmigrate.Config{})),
		},
	))

	options := newscripts.NewOptions(scriptsInRoot)
	nextScriptCmd := newscripts.NextScriptCmd(&newscripts.Config{
		Migration: migration,
		Options:   options,
		DB:        db,
		Objects: []any{
			sample(&models.UserV1{}, &models.UserV2{}, &models.UserV3{}),
			sample(&models.InfoV1{}, &models.InfoV2{}, &models.InfoV3{}),
		},
	})

	newMigrateCmd := cobramigration.NewMigrateCmd(migration)

	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "main",
		Long:  "main",
	}
	rootCmd.AddCommand(nextScriptCmd)
	rootCmd.AddCommand(newMigrateCmd)

	must.Done(rootCmd.Execute())
}

func sample(objects ...interface{}) any {
	must.Have(objects)
	idx := rand.IntN(len(objects))
	return objects[idx]
}

type MysqlConfig struct {
	Source   string
	Username string
	Password string
}

func newGormDB(cfg *MysqlConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@%s", must.Nice(cfg.Username), must.Nice(cfg.Password), must.Nice(cfg.Source))

	db := rese.P1(gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info, //设置级别
				IgnoreRecordNotFoundError: true,        //有些场景就是需要查不到时才创建，因此查询不到时不算错误
				Colorful:                  true,
			}),
			TranslateError: true,
		},
	))
	sqlDB := rese.P1(db.DB())
	sqlDB.SetConnMaxIdleTime(60 * time.Second)
	sqlDB.SetMaxIdleConns(500)

	zaplog.SUG.Debugln("正在检查数据库连接")
	must.Done(sqlDB.Ping())
	zaplog.SUG.Debugln("已经检查数据库连接")
	return db
}
