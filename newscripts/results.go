package newscripts

import "github.com/go-xlan/go-migrate/checkmigration"

type ScriptAction string

const (
	CreateScript ScriptAction = "create-script"
	UpdateScript ScriptAction = "update-script"
)

type NextScript struct {
	Action      ScriptAction
	ForwardName string
	ReverseName string
}

func (next *NextScript) WriteScripts(migrationOps checkmigration.MigrationOps, options *Options) {
	forwardScript := migrationOps.GetForwardScript()
	mustWriteScript(next.Action, next.ForwardName, forwardScript, options)

	reverseScript, _ := migrationOps.GetReverseScript()
	mustWriteScript(next.Action, next.ReverseName, reverseScript, options)
}
