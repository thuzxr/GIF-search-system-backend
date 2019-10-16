package cache

import (
	// "fmt"
	"os"
	"io/ioutil"
	"strings"
	"encoding/base64"
	"backend/utils"
	"strconv"
	// jsoniter "github.com/json-iterator/go"
)

func FastWrite(filepath string,content []byte){
	w1, _:=os.OpenFile("cache",os.O_CREATE|os.O_TRUNC,0644)
	_, _=w1.Write(content)
	_=w1.Close()
}

func FastAppend(filepath string,content []byte){
	w1, _:=os.OpenFile("cache",os.O_APPEND,0644)
	_, _=w1.Write(content)
	_=w1.Close()
}

func OfflineCacheInit(){
	_,err:=os.Stat("cacheList")
	if os.IsNotExist(err){
		_=os.Mkdir("cacheList",os.ModePerm)
	}
}

func OfflineCacheAppend(keyword string,gif []utils.Gifs){
	w1, _:=os.OpenFile("cacheList/"+base64.URLEncoding.EncodeToString([]byte(keyword)),os.O_CREATE|os.O_TRUNC,0644)
	_, _=w1.Write([]byte(strconv.FormatInt(int64(len(gif)),10)+" "))
	for i:=0;i<len(gif);i++{
		_, _=w1.Write([]byte(gif[i].Name+" "))
	}
	_=w1.Close()
}

func OfflineCacheQuery(keyword string) []string{
	var res []string
	fname:="cacheList/"+base64.URLEncoding.EncodeToString([]byte(keyword))
	_,err:=os.Stat(fname)
	if os.IsNotExist(err){
		res=append(res,"Failed")
		return res
	}
	res=append(res,"Succeed")
	ind,_:=ioutil.ReadFile(fname)
	res=append(res,strings.Fields(string(ind))...)
	return res
}

// func main(){
// 	ind,_:=ioutil.ReadFile("cache")
// 	lis:=strings.Fields(string(ind))
// 	fmt.Println(lis)
// 	// fmt.Println(base64.URLEncoding.EncodeToString([]byte("哈哈")))
// 	var gif0 []utils.Gifs
// 	gif0=append(gif0,make([]utils.Gifs,3)...)
// 	gif0[0].Name="111"
// 	gif0[1].Name="222"
// 	gif0[2].Name="333"
// 	OfflineCacheInit()
// 	// OfflineCacheAppend("哈哈哈",gif0)
// 	fmt.Println(OfflineCacheQuery("哈哈哈"))
// }

