package main

import (
	"backend/cache"
	"backend/ossUpload"
	"backend/search"
	"backend/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	// "backend/word"
)

// func SearchDemo(searchKey string, gifs []utils.Gifs) string {
// 	for i := range gifs {
// 		for j := range gifs[i].Keyword {
// 			if strings.Compare(gifs[i].Keyword[j], searchKey) == 0 {
// 				return ossUpload.OssSignLink(gifs[i], 3600)
// 			}
// 		}
// 	}
// 	return "Null"
// }

func RouterSet() *gin.Engine {
	r := gin.Default()
	names, titles, keywords := search.FastIndexParse()
	m := cache.OfflineCacheReload()
	// gif := utils.JsonParse(".")
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
		msg := c.DefaultQuery("msg", "000")
		fmt.Println(msg)
		c.JSON(200, gin.H{
			"message": "hello world! --sent by GO",
		})
	})
	r.GET("/search", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
		// searchKey := c.Query("key")
		// res := SearchDemo(searchKey, gif)
		keyword := c.DefaultQuery("key", "UNK")
		// match := db.Query(DB, keyword)
		res, finded := m[keyword]
		var match []utils.Gifs
		if finded {
			match = res
			fmt.Println("Hit Cache " + keyword)
		} else {
			match = search.SimpleSearch(keyword, names, titles, keywords)
			go cache.OfflineCacheAppend(keyword, match)
		}
		for i := 0; i < len(match); i++ {
			match[i].Oss_url = ossUpload.OssSignLink(match[i], 3600)
		}
		if len(match) == 0 {
			c.JSON(200, gin.H{
				"status": "failed",
			})
		} else {
			c.JSON(200, gin.H{
				"status": "succeed",
				"result": match,
			})
		}

	})
	r.GET("/upload", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
		file := c.DefaultQuery("file", "defaultFile")
		fmt.Println(file)
		c.JSON(200, gin.H{
			"status": "succeed",
			"recept": file,
		})
	})
	// r.Run(":8000")
	return r
}

func main() {
	// gifs := utils.JsonParse(".")
	// DB := db.Connect_db()
	// db.CreateTable(DB)
	// db.DB_init(gifs, DB)

	// search.FastIndexInit()

	cache.OfflineCacheInit()
	r := RouterSet()
	r.Run(":80")

	// search.IndexInit(DB)
	// gifs:=search.IndexParse()
	// fmt.Println(gifs)
}
