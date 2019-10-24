package main

import (
	"backend/cache"
	"backend/ossUpload"
	"backend/recommend"
	"backend/search"
	"backend/upload"
	"backend/utils"
	"fmt"
	// "time"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"backend/word"
	"github.com/go-ego/gse"
)

func setHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
}

func RouterSet() *gin.Engine {
	cache.OfflineCacheInit()
	cache.OfflineCacheClear()
	r := gin.Default()
	gifs := utils.JsonParse("info.json")
	word2vec:=word.WordToVecInit()
	// name_reIdx:=word.Name_reIdx(gifs)
	// res:=word.WordToVec("静静地等红包",seg, m)
	re_idx, gif2vec, vec_h:=word.RankSearchInit()
	var seg gse.Segmenter
	seg.LoadDict()

	names, titles, keywords := search.FastIndexParse()
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

		// time0:=time.Now()
		keyword := c.DefaultQuery("key", "UNK")
		res, finded := m[keyword]
		var match []utils.Gifs
		// fmt.Println(time.Since(time0))
		if finded {
			match = res
			fmt.Println("Hit Cache " + keyword)
		} else {
			res:=word.RankSearch(keyword, word2vec, gif2vec, vec_h, re_idx, seg)
			// fmt.Println(time.Since(time0))
			match=make([]utils.Gifs,len(res))
			for i:=range(res){
				match[i]=maps[res[i]]
			}
			// fmt.Println(time.Since(time0))
			// match = search.SimpleSearch(keyword, names, titles, keywords)
			go cache.OfflineCacheAppend(keyword, match)
			// fmt.Println(time.Since(time0))
		}
		for i := 0; i < len(match); i++ {
			match[i].Oss_url = ossUpload.OssSignLink(match[i], 3600)
		}
		// fmt.Println(time.Since(time0))
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
	r := RouterSet()
	r.Run(":80")
}
