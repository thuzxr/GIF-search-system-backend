package cache

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastWriteAndAppend(t *testing.T) {
	filepath := "test_fast_func"
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

func TestOfflineCacheInit(t *testing.T) {
	gwd, _ := os.Getwd()
	os.Chdir(gwd + "..")
	OfflineCacheInit()
	assert.DirExists(t, "cache_name")
	assert.DirExists(t, "cache_title")
}
