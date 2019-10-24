package main

import (
	"backend/cache"
	"backend/ossUpload"
	"backend/recommend"
	"backend/search"
	"backend/upload"
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
	gifs := utils.JsonParse("info.json")
	fmt.Println(gifs[0])
	var maps map[string]utils.Gifs
	maps = make(map[string]utils.Gifs)
	for _, gif := range gifs {
		maps[gif.Name] = gif
	}

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
		keyword := c.DefaultQuery("keyword", "")
		name := c.DefaultQuery("name", "")
		title := c.DefaultQuery("title", "")
		keywords, names, titles = upload.Upload(keyword, name, title, keywords, names, titles)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/recommend", func(c *gin.Context) {
		name := c.DefaultQuery("name", "")
		recommend_gifs := recommend.Recommend(maps[name], gifs)
		for i := 0; i < len(recommend_gifs); i++ {
			recommend_gifs[i].Oss_url = ossUpload.OssSignLink(recommend_gifs[i], 3600)
		}
		c.JSON(200, gin.H{
			"status": "succeed",
			"result": recommend_gifs,
		})
	})

	return r
}

func main() {
	cache.OfflineCacheInit()
	r := RouterSet()
	r.Run(":8080")
}
