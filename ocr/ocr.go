package ocr

import (
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"github.com/json-iterator/go"
)

type Gifs struct {
	Name string
	Title string
	Keyword []string
	Gif_url string
	Cover_url string
	Oss_url string
}

//用于读取实例gif库的info.json，中期开发将替换为完整Gif库的链接，返回值是一个struct Gifs类
func JsonParse(path0 string) []Gifs{
	// path0 :="C:\\Users\\Mr Handsome\\Downloads\\Material_alt\\gif-dio-tmp\\gif_package"

	var gifs []Gifs

	bytes, _ := ioutil.ReadFile(path0+"\\info.json")

	json_data:=jsoniter.Get(bytes,"gifs")
	_data:=[]byte(json_data.ToString())

	size:=json_data.Size()

	gifs=append(gifs,make([]Gifs,size)...)

	for i:=0; i<size; i++{
		gifs[i].Name=jsoniter.Get(_data,i,"name").ToString()
		gifs[i].Title=jsoniter.Get(_data,i,"title").ToString()
		gifs[i].Keyword=strings.Fields(jsoniter.Get(_data,i,"keyword").ToString())
		gifs[i].Gif_url=jsoniter.Get(_data,i,"gif_url").ToString()
		gifs[i].Cover_url=jsoniter.Get(_data,i,"cover_url").ToString()
		gifs[i].Oss_url=""
	}

	return gifs
}

//适配OCR接口的函数，用于实现Gif图片的OCR结果，返回值是一个string切片
func Ocr(gif Gifs) []string{
	ocrSource:="http://lf.snssdk.com/2/wap/search/extra/ocr/"
	req, err:=http.NewRequest("GET",ocrSource,nil)
	if err!=nil {
		fmt.Print(err)
		panic(err)
	}

	q:=req.URL.Query()
	q.Add("url",gif.Gif_url)
	q.Add("token","toutiao_ocr")
	req.URL.RawQuery=q.Encode()
	// fmt.Println(req.URL.String())

	var resp *http.Response
	resp, err=http.DefaultClient.Do(req)
	if err!=nil{
		fmt.Print(err)
		panic(err)
	}
	defer resp.Body.Close()

	body,_err:=ioutil.ReadAll(resp.Body)
	if _err!=nil {
		fmt.Print(err)
		panic(err)
	}
	// fmt.Println(body)
	var tags []string
	
	json_data:=jsoniter.Get(body,"data","tags")
	_data:=[]byte(json_data.ToString())

	size:=json_data.Size()
	tags=append(tags,make([]string,size)...)

	for i:=0;i<size;i++{
		tags[i]=jsoniter.Get(_data,i,"tag").ToString()
	}

	return tags
}

// func main() {
// 	// q:=req.URL.Query()
// 	// q.add("")
// 	gifs:=jsonParse()
// 	fmt.Println(gifs[0])
// 	tags:=ocr(gifs[0])
// 	fmt.Println(tags)
// }