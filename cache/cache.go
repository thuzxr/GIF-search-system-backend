package cache

import (
	"backend/utils"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func cacheNamePath() string {
	gwd, _ := os.Getwd()
	return path.Join(gwd, utils.CACHE_DIR, "cache_name")
}

func cacheTitlePath() string {
	gwd, _ := os.Getwd()
	return path.Join(gwd, utils.CACHE_DIR, "cache_title")
}

//封装的快速文件读写，目前无用途
func FastWrite(filepath string, content []byte) {
	w1, _ := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	_, _ = w1.Write(content)
	_ = w1.Close()
}

func FastAppend(filepath string, content []byte) {
	w1, _ := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	_, _ = w1.Write(content)
	_ = w1.Close()
}

//初始化Cache存储目录
func OfflineCacheInit() {
	_, err := os.Stat(cacheNamePath())
	if os.IsNotExist(err) {
		_ = os.Mkdir(cacheNamePath(), os.ModePerm)
	}
	_, err = os.Stat(cacheTitlePath())
	if os.IsNotExist(err) {
		_ = os.Mkdir(cacheTitlePath(), os.ModePerm)
	}
}

//向Cache添加keyword及其搜索结果
func OfflineCacheAppend(keyword string, gif []utils.Gifs) {
	w1, _ := os.OpenFile(path.Join(cacheNamePath(), base64.URLEncoding.EncodeToString([]byte(keyword))),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	// _, _=w1.Write([]byte(strconv.FormatInt(int64(len(gif)),10)+"#"))
	for i := 0; i < len(gif); i++ {
		_, _ = w1.Write([]byte(gif[i].Name + "#"))
	}
	_ = w1.Close()
	w1, _ = os.OpenFile(path.Join(cacheTitlePath(), base64.URLEncoding.EncodeToString([]byte(keyword))),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	// _, _=w1.Write([]byte(strconv.FormatInt(int64(len(gif)),10)+"#"))
	for i := 0; i < len(gif); i++ {
		_, _ = w1.Write([]byte(gif[i].Title + "#"))
	}
	_ = w1.Close()
}

//查询keyword对应的Cache
func OfflineCacheQuery(keyword string) []string {
	var res []string
	fname := path.Join(cacheNamePath(), base64.URLEncoding.EncodeToString([]byte(keyword)))
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		res = append(res, "Failed")
		return res
	}
	res = append(res, "Succeed")
	ind, _ := ioutil.ReadFile(fname)
	res = append(res, strings.Split(string(ind), "#")...)
	return res
}

//读取目前已存储的Cache
func OfflineCacheReload() map[string][]utils.Gifs {
	// m_name:=make(map[string][]string)
	// m_title:=make(map[string][]string)
	m := make(map[string][]utils.Gifs)
	gif := new(utils.Gifs)
	var gifs []utils.Gifs
	// var res []string
	dir, _ := ioutil.ReadDir(cacheNamePath())
	for _, fi := range dir {
		if !fi.IsDir() {
			gifs = make([]utils.Gifs, 0)
			b, _ := base64.URLEncoding.DecodeString(fi.Name())
			b0, _ := ioutil.ReadFile(path.Join(cacheNamePath(), fi.Name()))
			lisName := strings.Split(string(b0), "#")
			b0, _ = ioutil.ReadFile(path.Join(cacheTitlePath(), fi.Name()))
			lisTitle := strings.Split(string(b0), "#")
			for i := 0; i < len(lisName); i++ {
				gif.Name = lisName[i]
				gif.Title = lisTitle[i]
				gifs = append(gifs, *gif)
			}
			m[string(b)] = gifs[0 : len(gifs)-1]
		}
	}
	return m
}

func OfflineCacheClear() {
	dir, _ := ioutil.ReadDir(cacheNamePath())
	var TmpName string
	for _, fi := range dir {
		if !fi.IsDir() {
			TmpName = fi.Name()
			os.Remove(path.Join(cacheNamePath(), TmpName))
			os.Remove(path.Join(cacheTitlePath(), TmpName))
		}
	}
}

func OfflineCacheDelete(keyword string) {
	KeywordName := base64.URLEncoding.EncodeToString([]byte(keyword))
	_, err := os.Stat(path.Join(cacheNamePath(), KeywordName))
	if os.IsNotExist(err) {
		return
	}
	err = os.Remove(path.Join(cacheNamePath(), KeywordName))
	if err != nil {
		fmt.Println(err)
	}
	os.Remove(path.Join(cacheTitlePath(), KeywordName))
}
