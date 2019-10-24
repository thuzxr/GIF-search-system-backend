package utils

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

type Gifs struct {
	Name      string
	Title     string
	Keyword   string
	Gif_url   string
	Cover_url string
	Oss_url   string
	Word_idx  [][]int32
}

const (
	USERNAME  = "wangziqi"
	PASSWORD  = "QWEasd123_"
	NETWORK   = "tcp"
	PORT      = "3306"
	SERVER    = "49.233.71.202"
	DATABASE  = "GIF_INFO"
	CACHE_DIR = "cache/"
	HAMMING_EDGE = 133
	HAMMING_DIV = 79
)

//用于读取实例gif库的info.json，中期开发将替换为完整Gif库的链接，返回值是一个struct Gifs类
func JsonParse(path0 string) []Gifs {
	var gifs []Gifs

	bytes, _ := ioutil.ReadFile(path0 + "/info.json")
	jsonData := jsoniter.Get(bytes, "gifs")
	_data := []byte(jsonData.ToString())

	size := jsonData.Size()
	gifs = append(gifs, make([]Gifs, size)...)
	for i := 0; i < size; i++ {
		gifs[i].Name = jsoniter.Get(_data, i, "name").ToString()
		gifs[i].Title = jsoniter.Get(_data, i, "title").ToString()
		gifs[i].Keyword = jsoniter.Get(_data, i, "keyword").ToString()
		gifs[i].Gif_url = jsoniter.Get(_data, i, "gif_url").ToString()
		gifs[i].Cover_url = jsoniter.Get(_data, i, "cover_url").ToString()
		gifs[i].Oss_url = ""
		gifs[i].Word_idx = nil
	}
	return gifs
}
