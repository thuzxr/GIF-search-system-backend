package search

import(
	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"backend/utils"
	"strings"
	// "fmt"
	jsoniter "github.com/json-iterator/go"
	// "backend/database"
	// "os"
	"io/ioutil"
)

type IndexStruct struct {
	Size int `json:"size"`
	Gifs []utils.Gifs `json:"gifs"`
}

//返回json格式的index存储的Gif列表
func IndexParse() []utils.Gifs{
	bytes, _ := ioutil.ReadFile("searchIndex.json")
	var indexStruct IndexStruct
	_=jsoniter.Unmarshal(bytes, &indexStruct)

	return indexStruct.Gifs
}

//读取.ind格式的index,下同
func NameIndex() []string{
	b,_ :=ioutil.ReadFile("ind_name.ind")
	return strings.Split(string(b),"#")
}

func TitleIndex() []string{
	b,_ :=ioutil.ReadFile("ind_title.ind")
	return strings.Split(string(b),"#")
}

func KeywordIndex() []string{
	b,_ :=ioutil.ReadFile("ind_keyword.ind")
	return strings.Split(string(b),"#")
}