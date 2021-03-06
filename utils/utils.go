package utils

import (
	"io/ioutil"
	"strconv"
	"strings"

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
	Recommend []int
}

const (
	USERNAME = "wangziqi"
	PASSWORD = "QWEasd123_"
	NETWORK  = "tcp"
	PORT     = "3306"
	// SERVER        = "49.233.71.202"
	DATABASE      = "GIF_INFO_NEW"
	CACHE_DIR     = "cache/"
	HAMMING_EDGE  = 133
	HAMMING_DIV   = 79
	COOKIE_EXPIRE = 3600
	COOKIE_DOMAIN = "www.gifxiv.com"
	COOKIE_SALT   = "The_World"
	// SSLHOST       = "49.233.71.202:8080"
	STATUS     = "status"
	SUCCEED    = "succeed"
	RESULT     = "result"
	FIRSTNAME  = "FirstName"
	LASTNAME   = "LastName"
	ZIPCODE    = "ZipCode"
	COUNTRY    = "Country"
	HEIGHT     = "Height"
	BIRTHDAY   = "Birthday"
	ScanFailed = "scan failed, err:%v\n"
	// COOKIE_DOMAIN = "183.173.58.166"
	// COOKIE_DOMAIN = "gif-dio-stardustcrusaders.app.secoder.net"
)

type readjson struct {
	name      string
	title     string
	keyword   string
	gif_url   string
	cover_url string
	recommend []int
}

type Profile struct {
	Email     string
	FirstName string
	LastName  string
	Addr      string
	ZipCode   string
	City      string
	Country   string
	About     string
	Height    string
	Birthday  string
}

type Like_based_sort struct {
	Gif  Gifs
	Like int
}

type LikeSlice []Like_based_sort

func (s LikeSlice) Len() int           { return len(s) }
func (s LikeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s LikeSlice) Less(i, j int) bool { return s[i].Like > s[j].Like }

//用于读取实例gif库的info.json，中期开发将替换为完整Gif库的链接，返回值是一个struct Gifs类
func JsonParse(path0 string) []Gifs {
	var gifs []Gifs

	bytes, _ := ioutil.ReadFile(path0)
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
		load_strings := strings.Split(jsoniter.Get(_data, i, "recommend").ToString(), " ")
		for _, s := range load_strings {
			load_num, _ := strconv.Atoi(s)
			gifs[i].Recommend = append(gifs[i].Recommend, load_num)
		}
	}
	return gifs
}
