package main

import (
	"backend/cache"
	"backend/database"
	"backend/management/login"
	"backend/management/register"
	"backend/ossUpload"
	"backend/recommend"
	"backend/search"
	"backend/upload"
	"backend/utils"
	"backend/cookie"
	"fmt"

	// "time"
	"backend/word"

	"github.com/gin-gonic/gin"
	"github.com/go-ego/gse"
	_ "github.com/go-sql-driver/mysql"
)

func setHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	// c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"));
	// c.Header("Access-Control-Allow-Credentials","true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
	c.Header("Access-Control-Expose-Headers", "Date, Set-Cookie")
}

func RouterSet() *gin.Engine {
	DB := database.ConnectDB()
	database.CreateTable(DB)

	cache.OfflineCacheInit()
	cache.OfflineCacheClear()
	r := gin.Default()
	gifs := utils.JsonParse("info.json")
	// AdSearch_Enabled := word.DataCheck()
	AdSearch_Enabled := false

	var gif2vec map[string][][]uint8
	var word2vec map[string][]uint8
	var re_idx []string
	var vec_h [][]uint64
	var seg gse.Segmenter

	if AdSearch_Enabled {
		fmt.Println("Advanced Searching Enabled")
		word2vec = word.WordToVecInit()
		re_idx, gif2vec, vec_h = word.RankSearchInit()
		seg.LoadDict()
	} else {
		fmt.Println("Index not found, Advanced Searching Disabled")
	}
	//names, titles, keywords := search.FastIndexParse()
	names:=make([]string,0)
	titles:=make([]string,0)
	keywords:=make([]string,0)

	fmt.Println(gifs[0])
	var maps map[string]utils.Gifs
	maps = make(map[string]utils.Gifs)
	for _, gif := range gifs {
		maps[gif.Name] = gif
	}

	m := cache.OfflineCacheReload()
	// gif := utils.JsonParse(".")

	goc := cookie.CookieCacheInit()

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
			if AdSearch_Enabled {
				res := word.RankSearch(keyword, word2vec, gif2vec, vec_h, re_idx, seg)
				match = make([]utils.Gifs, len(res))
				for i := range res {
					match[i] = maps[res[i]]
				}
				match = append(match, search.SimpleSearch(keyword, names, titles, keywords)...)
			} else {
				match = search.SimpleSearch(keyword, names, titles, keywords)
			}
			go cache.OfflineCacheAppend(keyword, match)
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
		
		// if(!cookie.CookieCheck(c.Request, goc)){
		// 	c.JSON(200, gin.H{
				
		// 	})
		// }

		keyword := c.DefaultQuery("keyword", "")
		name := c.DefaultQuery("name", "")
		title := c.DefaultQuery("title", "")
		keywords, names, titles = upload.Upload(keyword, name, title, keywords, names, titles)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/recommend", func(c *gin.Context) {
		setHeader(c)

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

	r.GET("/login", func(c *gin.Context) {
		setHeader(c)

		user := c.DefaultQuery("user", "")
		password := c.DefaultQuery("password", "")

		status := login.Login(user, password, DB)
		if(status=="登陆成功！"){
			c.SetCookie("user_name", user, 3600, "/", utils.COOKIE_DOMAIN,  false, false)
			cookie.CookieSet(user, goc)
		}else{
			c.SetCookie("user_name", "", 3600, "/", utils.COOKIE_DOMAIN, false, false)
		}
		c.JSON(200, gin.H{
			"status": status,
		})
	})

	r.GET("/register", func(c *gin.Context) {
		setHeader(c)

		status := register.Register(c, DB)
		c.JSON(200, gin.H{
			"status": status,
		})
	})

	r.GET("/write_cookie", func(c *gin.Context) {
		setHeader(c)

		c.SetCookie("user_cookie", "cookie0", 3600, "/", utils.COOKIE_DOMAIN, false, true)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/read_cookie", func(c *gin.Context) {
		setHeader(c)

		b:=cookie.CookieCheck(c.Request, goc)
		c.JSON(200, gin.H{
			"res": b,
		})
	})

	return r
}

func main() {
	cache.OfflineCacheInit()
	r := RouterSet()
	r.Run(":8080")
}
