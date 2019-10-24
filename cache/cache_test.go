package cache

import (
	"backend/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	name := path.Base(cacheNamePath())
	assert.Equal(t, name, "cache_name")
	title := path.Base(cacheTitlePath())
	assert.Equal(t, title, "cache_title")
}

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

func mockGif(name string) utils.Gifs {
	return utils.Gifs{name, "title", "keyword", "gifurl", "covurl", "ossurl", nil,nil}
}

func TestOfflineCache(t *testing.T) {
	// test init
	gwd, _ := os.Getwd()
	os.Chdir("..")
	OfflineCacheInit()
	assert.DirExists(t, path.Join(gwd, "cache_name"))
	assert.DirExists(t, path.Join(gwd, "cache_title"))

	// test append and query
	gifs := []utils.Gifs{mockGif("testGif1"), mockGif("testGif2")}
	OfflineCacheAppend("testGif", gifs)
	gifs = []utils.Gifs{mockGif("gif1"), mockGif("gif2")}
	OfflineCacheAppend("gif", gifs)
	res := OfflineCacheQuery("testGif")
	fmt.Println(len(res))
	assert.Equal(t, len(res), 4)
	assert.Equal(t, "Succeed", res[0])
	assert.Equal(t, "testGif1", res[1])
	assert.Equal(t, "testGif2", res[2])

	// test reload
	gifmap := OfflineCacheReload()
	assert.Equal(t, "gif1", gifmap["gif"][0].Name)
	assert.Equal(t, "gif2", gifmap["gif"][1].Name)
	assert.Equal(t, "testGif1", gifmap["testGif"][0].Name)
	assert.Equal(t, "testGif2", gifmap["testGif"][1].Name)
	assert.Equal(t, len(gifmap), 2)

	// test delete
	OfflineCacheDelete("testGif")
	gifmap = OfflineCacheReload()
	assert.Equal(t, len(gifmap), 1)
	assert.Equal(t, "gif1", gifmap["gif"][0].Name)
	assert.Equal(t, "gif2", gifmap["gif"][1].Name)

	// test clear
	OfflineCacheClear()
	gifmap = OfflineCacheReload()
	assert.Equal(t, len(gifmap), 0)
}
