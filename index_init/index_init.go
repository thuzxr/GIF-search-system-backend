package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"backend/utils"
	db "backend/database"
	"fmt"
	"os"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type IndexStruct struct {
	Size int `json:"size"`
	Gifs []utils.Gifs `json:"gifs"`
}
//生成json格式的index，由于速度问题暂时未采用
func IndexInit(DB *sql.DB){
	rows, err := DB.Query("Select Name,Title,Keyword from GIF_INFO")
	defer func() {
		if rows!=nil{
			rows.Close();
		}
	}()

	if err!=nil{
		panic(err);
	}

	gif := new(utils.Gifs)
	var ans []utils.Gifs

	for rows.Next() {
		if serr := rows.Scan(&gif.Name,&gif.Title,&gif.Keyword); serr != nil {
			fmt.Printf("scan failed, err:%v\n", serr)
			return
		}
		ans = append(ans, *gif)
	}

	var indexStruct IndexStruct
	indexStruct.Size=len(ans)
	indexStruct.Gifs=ans
	// fmt.Println(indexStruct)
	b, merr:=jsoniter.Marshal(indexStruct)
	if merr!=nil{
		fmt.Println(merr)
		panic(merr)
	}
	w1,_ := os.OpenFile("searchIndex.json",os.O_CREATE|os.O_TRUNC,0644)
	_,_= w1.Write(b)
	_=w1.Close()
	fmt.Println("Index Generated")
}

//生成.ind格式的Gif列表
func FastIndexInit(DB *sql.DB) {
	rows, err := DB.Query("Select Name,Title,Keyword from GIF_INFO")
	defer func() {
		if rows!=nil{
			rows.Close();
		}
	}()

	if err!=nil{
		panic(err);
	}

	var names []string
	var titles []string
	var keywords []string
	var name,title,keyword string

	for rows.Next() {
		if serr := rows.Scan(&name,&title,&keyword); serr != nil {
			fmt.Printf("scan failed, err:%v\n", serr)
			return
		}
		names=append(names,name)
		titles=append(titles,title)
		keywords=append(keywords,keyword)
	}

	w1, _:=os.OpenFile("ind_name.ind",os.O_CREATE|os.O_TRUNC,0644)
	_, _=w1.Write([]byte(strings.Join(names, "#")))
	_=w1.Close()
	// cache.FastWrite("ind_name.ind",[]byte(strings.Join(names, "#")))
	w1, _=os.OpenFile("ind_title.ind",os.O_CREATE|os.O_TRUNC,0644)
	_, _=w1.Write([]byte(strings.Join(titles, "#")))
	_=w1.Close()
	w1, _=os.OpenFile("ind_keyword.ind",os.O_CREATE|os.O_TRUNC,0644)
	_, _=w1.Write([]byte(strings.Join(keywords, "#")))
	_=w1.Close()
	fmt.Println("FastIndex Inited")
} 

func main(){
	DB:=db.Connect_db()
	FastIndexInit(DB)
}