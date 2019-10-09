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
}

const (
	USERNAME = "wangziqi"
	PASSWORD = "QWEasd123_"
	NETWORK  = "tcp"
	PORT     = "3306"
	SERVER   = "49.233.71.202"
	DATABASE = "GIF_INFO"
)

//用于读取实例gif库的info.json，中期开发将替换为完整Gif库的链接，返回值是一个struct Gifs类
func JsonParse(path0 string) []Gifs {
	// path0 :="C:\\Users\\Mr Handsome\\Downloads\\Material_alt\\gif-dio-tmp\\gif_package"

	var gifs []Gifs

	bytes, _ := ioutil.ReadFile(path0 + "/info.json")

	json_data := jsoniter.Get(bytes, "gifs")
	_data := []byte(json_data.ToString())

	size := json_data.Size()

	gifs = append(gifs, make([]Gifs, size)...)

	for i := 0; i < size; i++ {
		gifs[i].Name = jsoniter.Get(_data, i, "name").ToString()
		gifs[i].Title = jsoniter.Get(_data, i, "title").ToString()
		gifs[i].Keyword = jsoniter.Get(_data, i, "keyword").ToString()
		gifs[i].Gif_url = jsoniter.Get(_data, i, "gif_url").ToString()
		gifs[i].Cover_url = jsoniter.Get(_data, i, "cover_url").ToString()
		gifs[i].Oss_url = ""
	}
	return gifs
}
