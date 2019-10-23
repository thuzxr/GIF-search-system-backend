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

func setHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
}

func RouterSet() *gin.Engine {
	r := gin.Default()
	names, titles, keywords := search.FastIndexParse()
	m := cache.OfflineCacheReload()
	// gif := utils.JsonParse(".")
	r.GET("/", func(c *gin.Context) {
		setHeader(c)

		msg := c.DefaultQuery("msg", "000")
		fmt.Println(msg)
		c.JSON(200, gin.H{
			"message": "hello world! --sent by GO",
		})
	})
	r.GET("/search", func(c *gin.Context) {
		setHeader(c)

		keyword := c.DefaultQuery("key", "UNK")
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
		setHeader(c)

		file := c.DefaultQuery("file", "defaultFile")
		fmt.Println(file)
		c.JSON(200, gin.H{
			"status": "succeed",
			"recept": file,
		})
	})
	return r
}

func main() {
	cache.OfflineCacheInit()
	r := RouterSet()
	r.Run(":80")
	// gifs := utils.JsonParse("/Users/saberrrrrrrr/go/src/backend")
	// model := cbow.Init("/Users/saberrrrrrrr/go/src/backend/tensorflow/python_models/CBOW", "/Users/saberrrrrrrr/go/src/backend/tensorflow/python_models/data/word2idx.json", gifs)
	// fmt.Println("recomend")
	// commend := cbow.Recommend(gifs[0], gifs, model)
	// fmt.Println(gifs[0])
	// fmt.Println(commend)
}
