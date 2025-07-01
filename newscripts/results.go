package newscripts

import (
	"path/filepath"

	"github.com/go-xlan/go-migrate/checkmigration"
	"github.com/yyle88/osexistpath/osmustexist"
)

type ScriptAction string

const (
	CreateScript ScriptAction = "create-script"
	UpdateScript ScriptAction = "update-script"
)

type NextScriptInfo struct {
	Action      ScriptAction
	ForwardName string
	ReverseName string
}

func (scriptInfo *NextScriptInfo) WriteScripts(migrationOps checkmigration.MigrationOps, options *Options) {
	forwardScript := migrationOps.GetForwardScript()
	mustWriteScript(scriptInfo.Action, scriptInfo.ForwardName, forwardScript, options)

	reverseScript, _ := migrationOps.GetReverseScript()
	mustWriteScript(scriptInfo.Action, scriptInfo.ReverseName, reverseScript, options)
}

func (scriptInfo *NextScriptInfo) ScriptExists(options *Options) bool {
	if osmustexist.IsFile(filepath.Join(options.ScriptsInRoot, scriptInfo.ForwardName)) {
		return true
	}
	if osmustexist.IsFile(filepath.Join(options.ScriptsInRoot, scriptInfo.ReverseName)) {
		return true
	}
	return false
}

func (scriptInfo *NextScriptInfo) GetScriptNames() *NextScriptNames {
	return &NextScriptNames{
		ForwardName: scriptInfo.ForwardName,
		ReverseName: scriptInfo.ReverseName,
	}
}
