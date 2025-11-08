package main

import (
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/go-xlan/go-migrate/cobramigration"
	"github.com/go-xlan/go-migrate/internal/demos/demo2x/internal/models"
	"github.com/go-xlan/go-migrate/newmigrate"
	"github.com/go-xlan/go-migrate/newscripts"
	"github.com/go-xlan/go-migrate/previewmigrate"
	"github.com/golang-migrate/migrate/v4"
	postgresmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// docker run -d --name=postgres -e POSTGRES_PASSWORD=123456 -p 5432:5432 postgres
	cfg := &PostgresConfig{
		Dsn: "postgres://postgres:123456@localhost:5432/xlan_migrate_demo2x?sslmode=disable&TimeZone=UTC",
	}
	scriptsInRoot := runpath.PARENT.Join("scripts")

	// Lazy initialization: database connection created only when command runs
	// 延迟初始化：仅在命令运行时才创建数据库连接
	getDB := func() *gorm.DB {
		db := newGormDB(cfg)
		return db
	}

	// Migration factory accepts database connection to share single connection (avoiding duplicate connections)
	// 迁移工厂接受数据库连接以共享单个连接（避免重复连接）
	getMigration := func(db *gorm.DB) *migrate.Migrate {
		sqlDB := rese.P1(db.DB())
		migrationDriver := rese.V1(postgresmigrate.WithInstance(sqlDB, &postgresmigrate.Config{}))
		return rese.P1(newmigrate.NewWithScriptsAndDatabase(
			&newmigrate.ScriptsAndDatabaseParam{
				ScriptsInRoot:    scriptsInRoot,
				DatabaseName:     "postgres",
				DatabaseInstance: migrationDriver,
			},
		))
	}

	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "main",
		Long:  "main",
	}
	rootCmd.AddCommand(newscripts.NextScriptCmd(&newscripts.Config{
		GetMigration: getMigration,
		GetDB:        getDB,
		Options:      newscripts.NewOptions(scriptsInRoot),
		Objects: []any{
			randomSample(&models.UserV1{}, &models.UserV2{}, &models.UserV3{}),
			randomSample(&models.InfoV1{}, &models.InfoV2{}, &models.InfoV3{}),
		},
	}))
	rootCmd.AddCommand(cobramigration.NewMigrateCmd(getDB, getMigration))
	rootCmd.AddCommand(previewmigrate.NewPreviewCmd(getDB, getMigration, scriptsInRoot))

	must.Done(rootCmd.Execute())
}

func randomSample(objects ...interface{}) any {
	must.Have(objects)
	idx := rand.IntN(len(objects))
	return objects[idx]
}

type PostgresConfig struct {
	Dsn string
}

func newGormDB(cfg *PostgresConfig) *gorm.DB {
	db := rese.P1(gorm.Open(
		postgres.Open(cfg.Dsn),
		&gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
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
