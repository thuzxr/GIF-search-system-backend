package ocr

import (
	"backend/utils"
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

//适配OCR接口的函数，用于实现Gif图片的OCR结果，返回值是一个string切片
func Ocr(gif utils.Gifs) []string {
	ocrSource := "http://lf.snssdk.com/2/wap/search/extra/ocr/"
	req, err := http.NewRequest("GET", ocrSource, nil)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	q := req.URL.Query()
	q.Add("url", gif.Gif_url)
	q.Add("token", "toutiao_ocr")
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	defer resp.Body.Close()

	body, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		fmt.Print(err)
		panic(err)
	}

	var tags []string
	json_data := jsoniter.Get(body, "data", "tags")
	_data := []byte(json_data.ToString())

	size := json_data.Size()
	tags = append(tags, make([]string, size)...)

	for i := 0; i < size; i++ {
		tags[i] = jsoniter.Get(_data, i, "tag").ToString()
	}

	return tags
}
