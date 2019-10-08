package main

import (
	"github.com/gin-gonic/gin"
	"backend/ocr"
	// "fmt"
	"strings"
	"backend/ossUpload"
)

func SearchDemo(searchKey string,gifs []ocr.Gifs) string{
	for i:=range gifs{
		for j:=range gifs[i].Keyword{
			if strings.Compare(gifs[i].Keyword[j],searchKey)==0{
				return ossUpload.OssSignLink(gifs[i],3600)
			}
		}
	}
	return "Null"
}

func RouterSet() *gin.Engine{
	r:=gin.Default()
	gif:=ocr.JsonParse(".")
	r.GET("/", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
			c.JSON(200, gin.H{
				"message": "hello world! --sent by GO",
			})
		})
	r.GET("/search",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
		searchKey:=c.Query("key")
		res:=SearchDemo(searchKey,gif)
		if res=="Null"{
			c.JSON(200,gin.H{
				"status":"failed",
			})
		}else{
			c.JSON(200,gin.H{
				"status":"succeed",
				"result": res ,
			})
		}

	})
	// r.Run(":8000")
	return r
} 

func main() {
	r:=RouterSet()
	r.Run(":8000")

	// gifs:=ocr.JsonParse(".")
	// var gif []ocr.Gifs
	// fmt.Println(gif[0])
	// fmt.Println(SearchDemo("吐出来",gif))
}