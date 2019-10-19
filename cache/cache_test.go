package cache

import (
	"os"
	"path"
	"testing"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestFastWriteAndAppend(t *testing.T) {
	gwd, _ := os.Getwd()
	filepath := path.Join(gwd, "test_fast_cache")
	testStr := "test_fast"
	FastWrite(filepath, []byte(testStr))

	data, _ := ioutil.ReadFile(filepath)
	assert.Equal(t, testStr, string(data[:]))

	testStr2 := "test_append"
	FastAppend(filepath, []byte(testStr2))
	data, _ = ioutil.ReadFile(filepath)
	assert.Equal(t, testStr+testStr2, string(data[:]))
	os.Remove(filepath)
}

func TestOfflineCache(t *testing.T) {
	gwd, _ := os.Getwd()
	os.Chdir(path.Dir(gwd))
	OfflineCacheInit()
	assert.DirExists(t, path.Join(gwd, "cache_name"))
	assert.DirExists(t, path.Join(gwd, "cache_title"))
}
