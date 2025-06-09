package sketch1

import (
	"os"
	"testing"

	"github.com/go-xlan/go-migrate/internal/tests"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

func TestFileOpen(t *testing.T) {
	// cp from https://github.com/golang-migrate/migrate/blob/278833935c12dda022b1355f33a897d895501c45/source/file/file_test.go#L33
	f := file.File{}
	sourceInstance, err := f.Open("file://" + runpath.PARENT.Join("scripts"))
	require.NoError(t, err)
	defer rese.F0(sourceInstance.Close)

	version, err := sourceInstance.First()
	require.NoError(t, err)
	t.Log(version)

	for {
		t.Log(version)

		tests.ShowSourceContent(t, sourceInstance, version, source.Up)
		tests.ShowSourceContent(t, sourceInstance, version, source.Down)

		nextVersion, err := sourceInstance.Next(version)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return
			}
			t.FailNow()
		}
		version = nextVersion
	}
}
