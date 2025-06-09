package newscripts

import "github.com/go-xlan/go-migrate/checkmigration"

type ScriptAction string

const (
	CreateScript ScriptAction = "create-script"
	ModifyScript ScriptAction = "modify-script"
)

type ScriptNames struct {
	ForwardName string
	ReverseName string
}

type NextScript struct {
	ScriptAction ScriptAction
	Names        *ScriptNames
}

func (next *NextScript) WriteScripts(migrationOps checkmigration.MigrationOps, options *Options) {
	forwardScript := migrationOps.GetForwardScript()
	mustWriteScript(next.ScriptAction, next.Names.ForwardName, forwardScript, options)

	reverseScript, _ := migrationOps.GetReverseScript()
	mustWriteScript(next.ScriptAction, next.Names.ReverseName, reverseScript, options)
}
