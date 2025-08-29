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

// Config contains all necessary components for migration script generation via CLI
// Combines migration instance, options, database connection, and model objects
// Used to configure Cobra commands for script management operations
//
// Config 包含通过 CLI 进行迁移脚本生成所需的所有组件
// 结合迁移实例、选项、数据库连接和模型对象
// 用于配置 Cobra 命令进行脚本管理操作
type Config struct {
	Migration *migrate.Migrate // Migration instance for version control // 用于版本控制的迁移实例
	Options   *Options         // Script generation options // 脚本生成选项
	DB        *gorm.DB         // Database connection for schema analysis // 用于结构分析的数据库连接
	Objects   []interface{}    // GORM model objects for migration analysis // 用于迁移分析的 GORM 模型对象
}

// NextScriptCmd creates the main command for migration script management with subcommands
// Provides root command that displays current migration status and script information
// Includes create and update subcommands for comprehensive script management
//
// NextScriptCmd 创建带有子命令的迁移脚本管理主命令
// 提供显示当前迁移状态和脚本信息的根命令
// 包含用于全面脚本管理的创建和更新子命令
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

			scriptInfo := GetNextScriptInfo(config.Migration, config.Options, NewScriptNaming())
			zaplog.SUG.Infoln("next-script-info:", neatjsons.S(scriptInfo))

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

// createNewScriptCmd creates command for generating new migration scripts with version control
// Supports multiple version patterns (NEXT, UNIX, TIME) and custom descriptions
// Validates that new scripts should be created rather than updated
//
// createNewScriptCmd 创建用于生成带版本控制的新迁移脚本的命令
// 支持多种版本模式（NEXT、UNIX、TIME）和自定义描述
// 验证应该创建新脚本而不是更新现有脚本
func createNewScriptCmd(config *Config) *cobra.Command {
	var versionTypeInput string
	var description string
	var allowEmptyScript bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new migration script",
		Run: func(cmd *cobra.Command, args []string) {
			// 将字符串转换为 VersionPattern 枚举
			versionType := parseVersionType(versionTypeInput)

			// 创建 ScriptNaming（传入参数）
			scriptNaming := &ScriptNaming{
				VersionType: versionType,
				Description: description,
			}
			zaplog.SUG.Infoln("script-naming:", neatjsons.S(scriptNaming))

			// 获取下一组脚本名
			scriptInfo := GetNextScriptInfo(config.Migration, config.Options, scriptNaming)
			zaplog.SUG.Infoln("script-names:", neatjsons.S(scriptInfo.GetScriptNames()))

			// 假设系统建议你更新最新的脚本内容，而你选择的是创建，就报错
			if scriptInfo.Action == UpdateScript {
				eroticgo.RED.ShowMessage("FAILED. Use [update script] when THERE ARE UNMIGRATED SCRIPTS.")
				zaplog.SUG.Infoln(eroticgo.RED.Sprint("FAILED"))
				return
			}
			// 需要符合预期-避免出现其它情况，比如既非创建也非更新的其它情况
			must.Same(scriptInfo.Action, CreateScript)

			// 获取迁移操作并生成文件
			migrateOps := checkmigration.GetMigrateOps(config.DB, config.Objects)
			if len(migrateOps) > 0 || allowEmptyScript || scriptInfo.ScriptExists(config.Options) {
				scriptInfo.WriteScripts(migrateOps, config.Options)
			}

			eroticgo.GREEN.ShowMessage("SUCCESS")
		},
	}

	// 增加 flag 参数
	cmd.Flags().StringVarP(&versionTypeInput, "version-type", "t", "NEXT", "version pattern: NEXT, UNIX, TIME")
	cmd.Flags().StringVarP(&description, "description", "d", "script", "description for migration file name")
	cmd.Flags().BoolVarP(&allowEmptyScript, "allow-empty-script", "e", false, "allow creating script when no schema changes")

	return cmd
}

// updateTopScriptCmd creates command for updating the latest uncommitted migration script
// Updates existing script files with current database schema differences
// Validates that scripts exist and should be updated rather than newly created
//
// updateTopScriptCmd 创建用于更新最新未提交迁移脚本的命令
// 使用当前数据库结构差异更新现有脚本文件
// 验证脚本存在并应该被更新而不是新创建
func updateTopScriptCmd(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update top migration script",
		Run: func(cmd *cobra.Command, args []string) {
			scriptInfo := GetNextScriptInfo(config.Migration, config.Options, NewScriptNaming())
			zaplog.SUG.Infoln("script-names:", neatjsons.S(scriptInfo.GetScriptNames()))

			// 假设系统建议你创建最脚本内容，而你选择的是更新旧文件，就报错
			if scriptInfo.Action == CreateScript {
				eroticgo.RED.ShowMessage("FAILED. Use [create script] when THERE ARE NO UNMIGRATED SCRIPTS.")
				zaplog.SUG.Infoln(eroticgo.RED.Sprint("FAILED"))
				return
			}
			// 需要符合预期-避免出现其它情况，比如既非创建也非更新的其它情况
			must.Same(scriptInfo.Action, UpdateScript)

			migrateOps := checkmigration.GetMigrateOps(config.DB, config.Objects)
			if len(migrateOps) > 0 || scriptInfo.ScriptExists(config.Options) {
				scriptInfo.WriteScripts(migrateOps, config.Options)
			}
			eroticgo.GREEN.ShowMessage("SUCCESS")
		},
	}
}
