package newscripts

import "fmt"

type Options struct {
	ScriptsInRoot   string
	NewScriptPrefix func(version uint) (prefix string)
	DryRun          bool
	SurveyWritten   bool
	DefaultSuffix   string
}

func NewOptions(scriptsInRoot string) *Options {
	return &Options{
		ScriptsInRoot: scriptsInRoot,
		NewScriptPrefix: func(version uint) (prefix string) {
			return fmt.Sprintf("%05d_script", version)
		},
		DryRun:        false,
		SurveyWritten: false,
		DefaultSuffix: "sql",
	}
}
