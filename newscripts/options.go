package newscripts

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yyle88/must"
)

// Options contains configuration parameters for migration script generation and execution
// Controls script location, execution modes, and user interaction behavior
// Provides flexible configuration for different deployment and development scenarios
//
// Options 包含迁移脚本生成和执行的配置参数
// 控制脚本位置、执行模式和用户交互行为
// 为不同的部署和开发场景提供灵活配置
type Options struct {
	ScriptsInRoot string // Path to migration scripts DIR // 迁移脚本 DIR 路径
	DryRun        bool   // Enable dry-run mode without file writes // 启用试运行模式，不写入文件
	SurveyWritten bool   // Enable interactive confirmation prompts // 启用交互式确认提示
	DefaultSuffix string // Default file extension for scripts // 脚本的默认文件扩展名
}

// NewOptions creates default configuration for script generation with specified root DIR
// Sets up reasonable defaults for safe script generation and execution
// Returns configured options ready for customization
//
// NewOptions 为指定的根 DIR 创建脚本生成的默认配置
// 设置合理的默认值用于安全的脚本生成和执行
// 返回已配置的选项，可进一步自定义
func NewOptions(scriptsInRoot string) *Options {
	return &Options{
		ScriptsInRoot: scriptsInRoot,
		DryRun:        false,
		SurveyWritten: false,
		DefaultSuffix: "sql",
	}
}

// VersionPattern defines the strategy for generating migration script version numbers
// Supports different versioning approaches for various project requirements
// Each pattern provides unique benefits for different development workflows
//
// VersionPattern 定义生成迁移脚本版本号的策略
// 为各种项目需求支持不同的版本控制方法
// 每种模式为不同的开发工作流提供独特优势
type VersionPattern string

const (
	VersionNext VersionPattern = "NEXT" // Auto-incrementing numbers, e.g., 00001, 00002 // 自动递增编号，例如 00001, 00002
	VersionUnix VersionPattern = "UNIX" // Unix timestamp versions, e.g., 1678693920 // Unix 时间戳版本，例如 1678693920
	VersionTime VersionPattern = "TIME" // Formatted datetime versions, e.g., 20250621103045 // 格式化日时版本，例如 20250621103045
)

func parseVersionType(s string) VersionPattern {
	switch strings.ToUpper(s) {
	case "NEXT":
		return VersionNext
	case "UNIX":
		return VersionUnix
	case "TIME": // 支持 TIME 和 DATETIME 两种写法
		return VersionTime
	default:
		panic("unknown version-type: " + s + " (must be NEXT, UNIX, TIME)")
	}
}

// ScriptNaming contains configuration for migration script naming conventions
// Combines version generation strategy with descriptive naming
// Used to create consistent and meaningful script file names
//
// ScriptNaming 包含迁移脚本命名约定的配置
// 将版本生成策略与描述性命名相结合
// 用于创建一致且有意义的脚本文件名
type ScriptNaming struct {
	VersionType VersionPattern // Version number generation strategy // 版本号生成策略
	Description string         // Descriptive name for migration scripts // 迁移脚本的描述性名称
}

// NewScriptNaming creates default script naming configuration with incremental versioning
// Sets up standard naming patterns suitable for most migration scenarios
// Returns naming configuration ready for immediate use or further customization
//
// NewScriptNaming 创建具有增量版本控制的默认脚本命名配置
// 设置适用于大多数迁移场景的标准命名模式
// 返回可立即使用或进一步自定义的命名配置
func NewScriptNaming() *ScriptNaming {
	return &ScriptNaming{
		VersionType: VersionNext,
		Description: "script",
	}
}

func (T *ScriptNaming) newVersion(versionNum uint) string {
	switch T.VersionType {
	case VersionNext:
		return fmt.Sprintf("%05d", versionNum) // 在数字左侧补零宽度 5 位，例如 00001, 00002
	case VersionUnix:
		return strconv.FormatInt(time.Now().Unix(), 10) // 当前 Unix 时间戳（秒数）
	case VersionTime:
		return time.Now().Format("20060102150405") // 格式：YYYYMMDDHHMMSS，例如 20250621103045
	default:
		panic("unknown VersionPattern: " + string(T.VersionType)) // 如果没有匹配，返回空字符串或 panic
	}
}

func (T *ScriptNaming) NewScriptPrefix(version uint) string {
	return fmt.Sprintf("%s_%s", must.Nice(T.newVersion(version)), must.Nice(T.Description))
}
