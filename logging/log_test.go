package logging_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"iot/logging"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	defer os.RemoveAll(logging.LogPath)

	require.NoError(t, logging.Setup())

	testContent := "TEST_OUTPUT"
	log.Println(testContent)

	files, err := ioutil.ReadDir(logging.LogPath)
	require.NoError(t, err)
	assert.Equal(t, 1, len(files))
	for _, f := range files {
		assert.Equal(t, time.Now().Format("20060102"), f.Name())
		c, err := ioutil.ReadFile(path.Join(logging.LogPath, f.Name()))
		require.NoError(t, err)
		require.Contains(t, string(c), testContent)
	}
}