package main

import (
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/go-xlan/go-migrate/cobramigration"
	"github.com/go-xlan/go-migrate/internal/demos/demo1x/internal/models"
	"github.com/go-xlan/go-migrate/migrationparam"
	"github.com/go-xlan/go-migrate/migrationstate"
	"github.com/go-xlan/go-migrate/newmigrate"
	"github.com/go-xlan/go-migrate/newscripts"
	"github.com/go-xlan/go-migrate/previewmigrate"
	"github.com/golang-migrate/migrate/v4"
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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "main",
		Long:  "main",
	}

	cfg := &MysqlConfig{
		Dsn: "root:123456@tcp(localhost:3306)/xlan_migrate_demo1x?charset=utf8mb4&parseTime=true&multiStatements=true",
	}
	scriptsInRoot := runpath.PARENT.Join("scripts")

	migrationparam.SetDebugMode(true)

	// Migration connection with lazy initialization and unified resource management
	// 迁移连接，支持延迟初始化和统一资源管理
	param := migrationparam.NewMigrationParam(
		func() *gorm.DB {
			return newGormDB(cfg)
		},
		func(db *gorm.DB) *migrate.Migrate {
			sqlDB := rese.P1(db.DB())
			migrationDriver := rese.V1(mysqlmigrate.WithInstance(sqlDB, &mysqlmigrate.Config{}))
			return rese.P1(newmigrate.NewWithScriptsAndDatabase(
				&newmigrate.ScriptsAndDatabaseParam{
					ScriptsInRoot:    scriptsInRoot,
					DatabaseName:     "mysql",
					DatabaseInstance: migrationDriver,
				},
			))
		},
	)

	// Random version objects to simulate different development stages
	// 随机版本对象，模拟不同开发阶段的迁移场景
	objects := []any{
		randomSample(&models.UserV1{}, &models.UserV2{}, &models.UserV3{}),
		randomSample(&models.InfoV1{}, &models.InfoV2{}, &models.InfoV3{}),
	}

	rootCmd.AddCommand(newscripts.NewScriptCmd(&newscripts.Config{
		Param:   param,
		Options: newscripts.NewOptions(scriptsInRoot),
		Objects: objects,
	}))
	rootCmd.AddCommand(cobramigration.NewMigrateCmd(param))
	rootCmd.AddCommand(previewmigrate.NewPreviewCmd(param, scriptsInRoot))
	rootCmd.AddCommand(migrationstate.NewStatusCmd(&migrationstate.Config{
		Param:       param,
		ScriptsPath: scriptsInRoot,
		Objects:     objects,
	}))

	must.Done(rootCmd.Execute())
}

func randomSample(objects ...interface{}) any {
	must.Have(objects)
	idx := rand.IntN(len(objects))
	return objects[idx]
}

type MysqlConfig struct {
	Dsn string
}

func newGormDB(cfg *MysqlConfig) *gorm.DB {
	db := rese.P1(gorm.Open(
		mysql.Open(cfg.Dsn),
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
