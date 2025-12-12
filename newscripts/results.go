package newscripts

import (
	"path/filepath"

	"github.com/go-xlan/go-migrate/checkmigration"
	"github.com/yyle88/osexistpath/osmustexist"
)

// ScriptAction represents the type of action to perform on migration scripts
// Determines whether to create new scripts or update existing ones
//
// ScriptAction 表示对迁移脚本执行的操作类型
// 决定是创建新脚本还是更新现有脚本
type ScriptAction string

const (
	CreateScript ScriptAction = "create-script" // Create new migration scripts // 创建新的迁移脚本
	UpdateScript ScriptAction = "update-script" // Update existing scripts // 更新现有脚本
)

// NewScriptInfo contains information about the next migration script to be generated
// Includes action type and both forward and reverse script filenames
// Used to coordinate script creation and updates
//
// NewScriptInfo 包含将要生成的下一个迁移脚本的信息
// 包含操作类型和正向、反向脚本文件名
// 用于协调脚本创建和更新
type NewScriptInfo struct {
	Action      ScriptAction // Type of script action to perform // 要执行的脚本操作类型
	ForwardName string       // Filename for forward migration script // 正向迁移脚本的文件名
	ReverseName string       // Filename for reverse migration script // 反向迁移脚本的文件名
}

// WriteScripts generates and writes both forward and reverse migration scripts to file system
// Creates script content from migration operations and handles file writing
// Supports both create and update scenarios based on script action
//
// WriteScripts 生成并将正向和反向迁移脚本写入文件系统
// 从迁移操作创建脚本内容并处理文件写入
// 基于脚本操作支持创建和更新场景
func (scriptInfo *NewScriptInfo) WriteScripts(migrationOps checkmigration.MigrationOps, options *Options) {
	forwardScript := migrationOps.GetForwardScript()
	mustWriteScript(scriptInfo.Action, scriptInfo.ForwardName, forwardScript, options)

	reverseScript, _ := migrationOps.GetReverseScript()
	mustWriteScript(scriptInfo.Action, scriptInfo.ReverseName, reverseScript, options)
}

// ScriptExists checks whether migration script files already exist in the target DIR
// Verifies existence of both forward and reverse script files
// Returns true if any of the script files are found
//
// ScriptExists 检查迁移脚本文件是否已存在于目标 DIR 中
// 验证正向和反向脚本文件的存在性
// 如果找到任何脚本文件则返回 true
func (scriptInfo *NewScriptInfo) ScriptExists(options *Options) bool {
	if osmustexist.IsFile(filepath.Join(options.ScriptsInRoot, scriptInfo.ForwardName)) {
		return true
	}
	if osmustexist.IsFile(filepath.Join(options.ScriptsInRoot, scriptInfo.ReverseName)) {
		return true
	}
	return false
}

func (scriptInfo *NewScriptInfo) GetScriptNames() *NewScriptNames {
	return &NewScriptNames{
		ForwardName: scriptInfo.ForwardName,
		ReverseName: scriptInfo.ReverseName,
	}
}
