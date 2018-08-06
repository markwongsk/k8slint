package importalias_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/markwongsk/go-k8slint/k8slint/importalias"
	"github.com/stretchr/testify/require"
)

const (
	testdataDir = "testdata"
	errorFile   = "errors.txt"
	srcPath     = "src"
)

func TestImportaliasNoError(t *testing.T) {
	files, err := ioutil.ReadDir(testdataDir)
	require.NoError(t, err)

	for _, file := range files {
		tc := path.Join(testdataDir, file.Name())
		src := path.Join(tc, srcPath)
		fname := path.Join(tc, errorFile)
		f, err2 := ioutil.ReadFile(fname)
		require.NoError(t, err2)
		buf := &bytes.Buffer{}
		err = importalias.Run([]string{src}, false, buf)
		if err == nil || err.Error() == "" {
			require.Equal(t, string(f), buf.String())
		} else {
			require.Fail(t, err.Error())
		}
	}
}
