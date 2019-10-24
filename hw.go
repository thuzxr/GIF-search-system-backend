package main

import (
	"backend/cache"
	"backend/ossUpload"
	// "backend/search"
	// cbow "backend/tensorflow"
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
	gifs := utils.JsonParse(".")
	word2vec:=word.WordToVecInit()
	name_reIdx:=word.Name_reIdx(gifs)
	// res:=word.WordToVec("静静地等红包",seg, m)
	re_idx, gif2vec, vec_h:=word.RankSearchInit()
	var seg gse.Segmenter
	seg.LoadDict()

	// names, titles, keywords := search.FastIndexParse()
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
				match[i]=*name_reIdx[res[i]]
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
	r := RouterSet()
	r.Run(":8000")
	// gifs := utils.JsonParse(".")
	// model := cbow.Init("tensorflow/python_models/CBOW", "tensorflow/python_models/data/word2idx.json", gifs)
	// fmt.Println("recomend")
	// commend := cbow.Recommend(gifs[0], gifs, model)
	// m:=word.VecParse()
	// fmt.Println(m[gifs[0].Name])
	// m:=word.FastVecParse()
	// fmt.Println(word.HammingCode(m["ff0c1056353070f84f7dd126a335cf57"][0]))
	// fmt.Println(gif2vec["13a29d7df01e08d84f5ed3690db02b72"])
	// fmt.Println(res)
	// fmt.Println(gifs[0])
	// fmt.Println(commend)
}
