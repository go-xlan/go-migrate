package newscripts

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yyle88/must"
)

type Options struct {
	ScriptsInRoot string
	DryRun        bool
	SurveyWritten bool
	DefaultSuffix string
}

func NewOptions(scriptsInRoot string) *Options {
	return &Options{
		ScriptsInRoot: scriptsInRoot,
		DryRun:        false,
		SurveyWritten: false,
		DefaultSuffix: "sql",
	}
}

type VersionPattern string

const (
	VersionNext VersionPattern = "NEXT" // 自动递增编号，例如 00001, 00002
	VersionUnix VersionPattern = "UNIX" // 以 Unix 时间戳为版本号，例如 1678693920
	VersionTime VersionPattern = "TIME" // 格式化时间，年月日+时间，例如 202506211030
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

type ScriptNaming struct {
	VersionType VersionPattern
	Description string // {version}_{description}.up.sql {version}_{description}.down.sql 里的 description
}

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
		return time.Now().Format("200601021504") // 格式：YYYYMMDDHHMM，例如 202506211030
	default:
		panic("unknown VersionPattern: " + string(T.VersionType)) // 如果没有匹配，返回空字符串或 panic
	}
}

func (T *ScriptNaming) NewScriptPrefix(version uint) string {
	return fmt.Sprintf("%s_%s", must.Nice(T.newVersion(version)), must.Nice(T.Description))
}
