package utils

import (
	"encoding/hex"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/zaplog"
)

func NewUUID32s() string {
	u := uuid.New()
	return hex.EncodeToString(u[:])
}

func WhistleCause(cause error) {
	if cause != nil {
		if errors.Is(cause, migrate.ErrNoChange) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("NO MIGRATION FILES TO RUN"))
		} else if errors.Is(cause, migrate.ErrNilVersion) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("NO VERSION IN VERSION-TABLE(schema_migrations)"))
		} else if errors.Is(cause, os.ErrNotExist) {
			zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("MIGRATION FILES NOT FOUND"))
		} else {
			zaplog.SUG.Panicln(eroticgo.RED.Sprint("MIGRATION FAILED:"), cause)
		}
		return
	}
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint("MIGRATION SUCCESS"))
}
