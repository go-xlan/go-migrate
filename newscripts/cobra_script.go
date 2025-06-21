package newscripts

import (
	"github.com/go-xlan/go-migrate/checkmigration"
	"github.com/go-xlan/go-migrate/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"gorm.io/gorm"
)

type Config struct {
	Migration *migrate.Migrate
	Options   *Options
	DB        *gorm.DB
	Objects   []interface{} // object list of your gorm models
}

func NextScriptCmd(config *Config) *cobra.Command {
	// Create root command
	var rootCmd = &cobra.Command{
		Use:   "next-script",
		Short: "Create next migration script",
		Long:  "Create next migration script",
		Run: func(cmd *cobra.Command, args []string) {
			version, dirtyFlag, err := config.Migration.Version()
			utils.WhistleCause(err) //panic when cause is not expected
			if dirtyFlag {
				eroticgo.RED.ShowMessage(version, "(DIRTY)")
			} else {
				eroticgo.GREEN.ShowMessage(version)
			}

			nextScript := GetNextScriptName(config.Migration, config.Options, NewScriptNaming())
			zaplog.SUG.Infoln("next-script-name:", neatjsons.S(nextScript))

			migrationOps := checkmigration.GetMigrateOps(config.DB, config.Objects)
			if len(migrationOps) > 0 {
				if forwardScript := migrationOps.GetForwardScript(); true {
					zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(forwardScript))
				}
				if reverseScript, ok := migrationOps.GetReverseScript(); ok {
					zaplog.SUG.Debugln(eroticgo.AMBER.Sprint(reverseScript))
				}
			}
			eroticgo.GREEN.ShowMessage("SUCCESS")
		},
	}

	rootCmd.AddCommand(createNewScriptCmd(config)) // Add `create` command
	rootCmd.AddCommand(updateTopScriptCmd(config)) // Add `update` command

	return rootCmd
}

func createNewScriptCmd(config *Config) *cobra.Command {
	var versionTypeStr string
	var description string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new migration script",
		Run: func(cmd *cobra.Command, args []string) {
			// 将字符串转换为 VersionPattern 枚举
			versionType := parseVersionType(versionTypeStr)

			// 创建 ScriptNaming（传入参数）
			scriptNaming := &ScriptNaming{
				VersionType: versionType,
				Description: description,
			}

			// 获取下一组脚本名
			nextScript := GetNextScriptName(config.Migration, config.Options, scriptNaming)
			zaplog.SUG.Infoln("create-script:", neatjsons.S(nextScript))
			must.Same(nextScript.Action, CreateScript)

			// 获取迁移操作并生成文件
			migrateOps := checkmigration.GetMigrateOps(config.DB, config.Objects)
			if len(migrateOps) > 0 {
				nextScript.WriteScripts(migrateOps, config.Options)
			}

			eroticgo.GREEN.ShowMessage("SUCCESS")
		},
	}

	// 增加 flag 参数
	cmd.Flags().StringVarP(&versionTypeStr, "version-type", "t", "NEXT", "version pattern: NEXT, UNIX, TIME")
	cmd.Flags().StringVarP(&description, "description", "d", "script", "description for migration file name")

	return cmd
}

func updateTopScriptCmd(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update top migration script",
		Run: func(cmd *cobra.Command, args []string) {
			nextScript := GetNextScriptName(config.Migration, config.Options, NewScriptNaming())
			zaplog.SUG.Infoln("update-script:", neatjsons.S(nextScript))
			must.Same(nextScript.Action, UpdateScript) //需要符合预期

			migrateOps := checkmigration.GetMigrateOps(config.DB, config.Objects)
			if len(migrateOps) > 0 {
				nextScript.WriteScripts(migrateOps, config.Options)
			}
			eroticgo.GREEN.ShowMessage("SUCCESS")
		},
	}
}
