// Package newscripts: Intelligent migration script generation and management system
// Provides automated script creation with versioning and naming conventions
// Features smart version progression and content generation based on database schema changes
// Integrates with GORM struct analysis to generate appropriate migration scripts
//
// newscripts: 智能迁移脚本生成和管理系统
// 提供自动化脚本创建，具有版本控制和命名约定
// 具有基于数据库结构变化的智能版本进展和内容生成
// 与 GORM 模型分析集成，生成适当的迁移脚本
package newscripts

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/pkg/errors"
	"github.com/yyle88/done"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/must/mustnum"
	"github.com/yyle88/must/muststrings"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// enumMigrateState represents the current migration state of the database
// Used to determine appropriate version progression and script creation approach
//
// enumMigrateState 表示数据库的当前迁移状态
// 用于确定适当的版本进展和脚本创建策略
type enumMigrateState string

const (
	noneMigrated enumMigrateState = "none-migrated" // No migrations executed // 无先前迁移
	onceMigrated enumMigrateState = "once-migrated" // Has migration records // 有迁移历史
)

// GetNewScriptInfo analyzes current migration state and determines next script information
// Examines existing migration files and database version to calculate appropriate next action
// Returns script naming details and action type for migration script creation
//
// GetNewScriptInfo 分析当前迁移状态并确定下一个脚本信息
// 检查现有迁移文件和数据库版本来计算适当的下一步操作
// 返回用于迁移脚本创建的脚本命名详情和操作类型
func GetNewScriptInfo(migration *migrate.Migrate, options *Options, naming *ScriptNaming) *NewScriptInfo {
	var migrateState enumMigrateState
	version, dirtyFlag, err := migration.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			must.Zero(version)
			migrateState = noneMigrated
		} else {
			panic(erero.Wro(err))
		}
	} else {
		must.False(dirtyFlag)
		migrateState = onceMigrated
	}
	must.Nice(migrateState)
	mustnum.Gte(version, 0)

	migrations := newMigrationsFromPath(options.ScriptsInRoot)

	nextVersion, nextAction := obtainNextVersion(migrateState, version, migrations, options)
	mustnum.Gt(nextVersion, version)
	scriptNames := obtainScriptNames(nextVersion, nextAction, options, migrations, naming)
	checkScriptName(scriptNames, version)

	zaplog.SUG.Debugln("next-action:", nextAction)
	zaplog.SUG.Debugln("script-name:", neatjsons.S(scriptNames))

	return &NewScriptInfo{
		Action:      nextAction,
		ForwardName: scriptNames.ForwardName,
		ReverseName: scriptNames.ReverseName,
	}
}

// newMigrationsFromPath scans DIR and builds migrations collection from script files
//
// newMigrationsFromPath 扫描 DIR 并从脚本文件构建迁移集合
func newMigrationsFromPath(scriptsInRoot string) *source.Migrations {
	migrations := source.NewMigrations()
	for _, e := range rese.V1(os.ReadDir(scriptsInRoot)) {
		if e.IsDir() {
			continue
		}
		migration := rese.P1(source.DefaultParse(e.Name()))
		zaplog.SUG.Debugln("append migration to migrations:", "version:", migration.Version, "direction:", migration.Direction)
		must.True(migrations.Append(migration))
	}
	return migrations
}

// mustWriteScript writes migration script to file system with validation and confirmation
// Validates file paths and handles both create and update scenarios
// Supports dry-run mode and interactive confirmation to ensure safe script generation
//
// mustWriteScript 将迁移脚本写入文件系统，具有安全检查和用户确认
// 验证文件路径并处理创建和更新场景
// 支持试运行模式和交互式确认，实现安全的脚本生成
func mustWriteScript(nextAction ScriptAction, shortName string, script string, options *Options) {
	var path = filepath.Join(options.ScriptsInRoot, shortName)
	if nextAction == CreateScript {
		must.False(osmustexist.IsFile(path))
	} else {
		must.Same(nextAction, UpdateScript)
		osmustexist.FILE(path)
	}
	zaplog.SUG.Debugln("path:", path, "script:", script)
	if options.DryRun {
		zaplog.SUG.Debugln("dry-run mode", options.DryRun)
		return
	}
	if options.SurveyWritten {
		var written bool
		prompt := &survey.Confirm{
			Message: "write script to path?",
			Default: true,
		}
		done.Done(survey.AskOne(prompt, &written))
		if !written {
			zaplog.SUG.Debugln("input_written", written)
			return
		}
	}
	// when file exist WriteFile truncates it before writing, without changing permissions.
	must.Done(os.WriteFile(path, []byte(script), 0644))
	zaplog.SUG.Debugln("done")
}

// checkScriptName validates script names match expected version sequence
//
// checkScriptName 验证脚本名称匹配预期的版本序列
func checkScriptName(scriptNames *NewScriptNames, previousVersion uint) {
	zaplog.LOG.Debug("check", zap.String("forward_name", scriptNames.ForwardName))
	mig1 := rese.P1(source.DefaultParse(must.Nice(scriptNames.ForwardName)))
	mustnum.Gt(mig1.Version, previousVersion)

	mig2 := rese.P1(source.DefaultParse(must.Nice(scriptNames.ReverseName)))
	mustnum.Gt(mig2.Version, previousVersion)

	must.Same(mig1.Version, mig2.Version)
}

// NewScriptNames holds forward and reverse script filenames
//
// NewScriptNames 持有正向和反向脚本文件名
type NewScriptNames struct {
	ForwardName string // Forward migration script filename // 正向迁移脚本文件名
	ReverseName string // Reverse migration script filename // 反向迁移脚本文件名
}

// obtainScriptNames generates script filenames based on action type and naming rules
//
// obtainScriptNames 基于操作类型和命名规则生成脚本文件名
func obtainScriptNames(nextVersion uint, nextAction ScriptAction, options *Options, migrations *source.Migrations, naming *ScriptNaming) *NewScriptNames {
	var scriptNames = &NewScriptNames{}
	switch nextAction {
	case CreateScript:
		prefix := naming.NewScriptPrefix(nextVersion)
		muststrings.Contains(prefix, "_")
		must.True(regexp.MustCompile(`^([0-9]+)_(.*)$`).MatchString(prefix))
		muststrings.NotContains(prefix, ".")

		// use first up-script file name suffix as new file name suffix
		suffix, ok := obtainFirstUpScriptNameSuffix(migrations)
		if !ok {
			suffix = options.DefaultSuffix
		}
		muststrings.NotContains(suffix, ".")

		scriptNames.ForwardName = fmt.Sprintf("%s.%s.%s", prefix, source.Up, suffix)
		must.True(source.DefaultRegex.MatchString(scriptNames.ForwardName))

		scriptNames.ReverseName = fmt.Sprintf("%s.%s.%s", prefix, source.Down, suffix)
		must.True(source.DefaultRegex.MatchString(scriptNames.ReverseName))
	case UpdateScript:
		scriptNames.ForwardName = resb.P1(migrations.Up(nextVersion)).Raw   // 123_name.up.ext
		scriptNames.ReverseName = resb.P1(migrations.Down(nextVersion)).Raw // 123_name.down.ext
	default:
		panic(erero.Errorf("IMPOSSIBLE case-value=%v", nextAction))
	}
	return scriptNames
}

// obtainFirstUpScriptNameSuffix extracts file extension from first migration script
//
// obtainFirstUpScriptNameSuffix 从第一个迁移脚本中提取文件扩展名
func obtainFirstUpScriptNameSuffix(migrations *source.Migrations) (string, bool) {
	firstVersion, ok := migrations.First()
	if ok {
		migration, ok := migrations.Up(firstVersion)
		if ok && migration != nil {
			matches := source.DefaultRegex.FindStringSubmatch(migration.Raw)
			if len(matches) == 5 {
				return matches[4], true
			}
		}
	}
	return "", false
}

// obtainNextVersion determines next version number and action based on migration state
//
// obtainNextVersion 基于迁移状态确定下一个版本号和操作
func obtainNextVersion(migrateState enumMigrateState, previousVersion uint, migrations *source.Migrations, options *Options) (uint, ScriptAction) {
	var nextVersion uint
	var ok bool
	switch migrateState {
	case noneMigrated:
		must.Zero(previousVersion)
		nextVersion, ok = migrations.First() // Get first script when no migrations exist // 假如从没做过就取首个脚本为待修改的
	case onceMigrated:
		nextVersion, ok = migrations.Next(previousVersion) // Get next version script // 否则就取下个版本的为待修改的
	default:
		panic(erero.Errorf("IMPOSSIBLE case-value=%v", migrateState))
	}
	if !ok {
		must.Zero(nextVersion)
		nextVersion = previousVersion + 1 // Return next version reference, can use timestamp etc. instead // 返回新版本号的参考值，当然后面也可以不使用这个参考值，而使用时间戳等版本号
		return nextVersion, CreateScript  // No script found, need to create new one // 假如取不到，就说明需要新建个脚本写内容
	}
	// if !options.ForceEdit {
	mustNoNextNextVersion(migrations, nextVersion) // Ensure this version is the latest, not intermediate // 需要确认获得的这个版本号就是最高的，而不是中间的，你也只能修改最高的
	// }
	return nextVersion, UpdateScript
}

// mustNoNextNextVersion ensures no versions exist after the given version
//
// mustNoNextNextVersion 确保给定版本之后不存在其他版本
func mustNoNextNextVersion(migrations *source.Migrations, nextVersion uint) {
	nextNextVersion, ok := migrations.Next(nextVersion)
	if !ok {
		return // Expected: no version after this means this is the latest // 这才是我们需要的，即没有下下个版本号的时候，就认为下个版本号就是最新的版本号
	}
	zaplog.LOG.Panic("script-is-not-latest-version", zap.Uint("next_version", nextVersion), zap.Uint("next_next_version", nextNextVersion))
}
