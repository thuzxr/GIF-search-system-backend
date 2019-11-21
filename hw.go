package main

import (
	"backend/cache"
	"strings"

	// "backend/cookie"
	"backend/cookie"
	"backend/database"
	"backend/management/login"
	"backend/management/register"
	"backend/management/vericode"
	"backend/ossUpload"
	"backend/recommend"
	"backend/search"

	// "backend/upload"
	"backend/utils"

	// "backend/cookie"
	"fmt"
	"time"

	"backend/word"

	"github.com/gin-gonic/gin"
	"github.com/go-ego/gse"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/dgrijalva/jwt-go/request"
)

func setHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	// c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"));
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
	// c.Header("Access-Control-Expose-Headers", "Date, set-cookie")
	// c.Header("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
}

func RouterSet() *gin.Engine {
	DB := database.ConnectDB()
	database.Init(DB)

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

	fmt.Println("OssUpdating")
	ossUpload.OssUpdate(gifs)
	fmt.Println("OssUpdated")
	fmt.Println(gifs[0].Oss_url)
	// ch_ossUpdate:=make(chan bool)

	go func() {
		for {
			time.Sleep(45 * time.Minute)
			fmt.Println("OssUpdating")
			ossUpload.OssUpdate(gifs)
			fmt.Println("OssUpdated")
		}
	}()

	if AdSearch_Enabled {
		fmt.Println("Advanced Searching Enabled")
		word2vec = word.WordToVecInit()
		re_idx, gif2vec, vec_h = word.RankSearchInit()
		seg.LoadDict()
	} else {
		fmt.Println("Index not found, Advanced Searching Disabled")
	}

	users, names, titles, _, keywords := database.LoadAll(DB)
	fmt.Println("total gifs size ", len(users))

	ch_gifUpdate := make(chan bool)
	go func() {
		for {
			select {
			case <-ch_gifUpdate:
				if true {
					users2, names2, titles2, infos2, keywords2 := database.LoadAll(DB)
					users = users2
					names = names2
					titles = titles2
					_ = infos2
					keywords = keywords2
				}
				fmt.Println("gif updates here")
				fmt.Println("total gifs size ", len(users))

				ch_gifUpdate <- false
				break
			default:
				break
			}
		}
	}()

	// names, titles, keywords := search.FastIndexParse()
	// names:=make([]string,0)
	// titles:=make([]string,0)
	// keywords:=make([]string,0)

	fmt.Println(gifs[0])
	var maps map[string]int
	maps = make(map[string]int)
	for i := range gifs {
		maps[gifs[i].Name] = i
	}

	m := cache.OfflineCacheReload()
	// gif := utils.JsonParse(".")

	// goc := cookie.CookieCacheInit()

	//Routers without Auth

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
					match[i] = gifs[maps[res[i]]]
				}
				// match = append(match, search.SimpleSearch(keyword, names, titles, keywords)...)
			} else {
				match0 := search.SimpleSearch(keyword, names, titles, keywords)
				match = make([]utils.Gifs, len(match0))
				for i := range match0 {
					match[i] = gifs[maps[match0[i].Name]]
				}
			}
			m[keyword] = match
			go cache.OfflineCacheAppend(keyword, match)
		}
		// for i := 0; i < len(match); i++ {
		// 	match[i].Oss_url = ossUpload.OssSignLink(match[i], 3600)
		// }
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

	r.POST("/login", cookie.UserAntiAuth(), func(c *gin.Context) {
		setHeader(c)

		user := c.DefaultPostForm("user", "")
		password := c.DefaultPostForm("password", "")

		status := login.Login(user, password, DB)

		if status != -1 {
			// c.SetCookie("user_name", string(cookie.ShaConvert(user)), 3600, "/", utils.COOKIE_DOMAIN,  false, false)
			// cookie.CookieSet(user, goc)
			cookie.TokenSet(c, user, status)
			favors := database.QueryFavor(user, DB)
			profile := database.QueryProfile(user, DB)

			c.JSON(200, gin.H{
				"status":    status,
				"Email":     profile[0],
				"FirstName": profile[1],
				"LastName":  profile[2],
				"Addr":      profile[3],
				"ZipCode":   profile[4],
				"City":      profile[5],
				"Country":   profile[6],
				"About":     profile[7],
				"Height":    profile[8],
				"Birthday":  profile[9],
				"favor":     favors,
			})
		} else {
			c.JSON(406, gin.H{
				"status": -1,
			})
			// c.SetCookie("user_name", "", 3600, "/", utils.COOKIE_DOMAIN, false, false)
		}
	})
	r.POST("/register", func(c *gin.Context) {
		setHeader(c)
		status := register.Register(c, DB)
		c.JSON(200, gin.H{
			"status": status,
		})
	})

	r.GET("/user_status", func(c *gin.Context) {
		res, _ := c.Request.Cookie("token")
		var status int
		var claims *cookie.MyClaims
		if res == nil {
			status = 0
			claims = &cookie.MyClaims{}
		} else {
			tokenString := res.Value
			fmt.Println(tokenString)
			claims, status = cookie.ClaimsParse(tokenString)
		}
		c.JSON(200, gin.H{
			"status": status,
			"claims": claims,
		})
	})

	r.GET("/user_login", func(c *gin.Context) {
		user := "user0"
		access := 1
		cookie.TokenSet(c, user, access)
		c.JSON(200, gin.H{
			"status": "success",
		})
	})

	r.GET("/refresh_veri", func(c *gin.Context) {
		setHeader(c)
		vericode.Get_vericode(c)
	})

	r.GET("/get_veri/:captchId", func(c *gin.Context) {
		setHeader(c)
		vericode.Gen_vericode(c)
	})

	//Routers with Auth

	r.Use(cookie.UserAuth())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"res": "test",
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		setHeader(c)
		c.SetCookie("token", "", -1, "/", utils.COOKIE_DOMAIN, false, false)
		c.JSON(200, gin.H{
			"status": "Logout",
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		setHeader(c)

		user := cookie.Getusername(c)
		info := c.DefaultPostForm("info", "")
		keyword := c.DefaultPostForm("keyword", "")
		name := c.DefaultPostForm("name", "")
		title := c.DefaultPostForm("title", "")
		// users, names, titles, infos, keywords = upload.Upload(users, names, titles, infos, keywords, user, name, title, info, keyword)
		// database.InsertGIF(DB, user, name, keyword, info, title)
		database.InsertUnderVerifyGIF(DB, user, name, keyword, info, title)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/toBeVerify", func(c *gin.Context) {
		setHeader(c)

		res := database.GetToVerifyGIF(DB)
		c.JSON(200, gin.H{
			"status": "succeed",
			"result": res,
		})
	})

	r.POST("/verify", func(c *gin.Context) {
		setHeader(c)
		name := c.DefaultPostForm("name", "")
		database.VerifyGIF(DB, name, ch_gifUpdate)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.POST("/remove", func(c *gin.Context) {
		setHeader(c)
		name := c.DefaultPostForm("name", "")
		database.DeleteGif(name, DB)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/recommend", func(c *gin.Context) {
		setHeader(c)

		name := c.DefaultQuery("name", "")
		recommend_gifs := recommend.Recommend(gifs[maps[name]], gifs)
		// for i := 0; i < len(recommend_gifs); i++ {
		// 	recommend_gifs[i].Oss_url = ossUpload.OssSignLink(recommend_gifs[i], 3600)
		// }
		c.JSON(200, gin.H{
			"status": "succeed",
			"result": recommend_gifs,
		})
	})

	r.GET("/profile", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		profile := database.QueryProfile(user, DB)
		c.JSON(200, gin.H{
			"Email":     profile[0],
			"FirstName": profile[1],
			"LastName":  profile[2],
			"Addr":      profile[3],
			"ZipCode":   profile[4],
			"City":      profile[5],
			"Country":   profile[6],
			"About":     profile[7],
			"Height":    profile[8],
			"Birthday":  profile[9],
		})
	})

	r.GET("/favor", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)

		favors := database.QueryFavor(user, DB)

		var osslinks []string
		for favor_id := range favors {
			favor := favors[favor_id]
			osslinks = append(osslinks, gifs[maps[favor]].Oss_url)
		}
		c.JSON(200, gin.H{
			"favors":  favors,
			"OssLink": osslinks,
		})
	})

	r.POST("/insert_favor", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		gifid := c.DefaultPostForm("GifId", "")
		favors := database.InsertFavor(user, gifid, DB)
		c.JSON(200, gin.H{
			"status": favors,
		})
	})

	r.POST("/delete_favor", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		gifid_string := c.DefaultPostForm("GifId", "")
		gifids := strings.Split(gifid_string, " ")
		favors := database.DeleteFavor(user, gifids, DB)
		c.JSON(200, gin.H{
			"status": favors,
		})
	})

	r.GET("/follow", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)

		follows := database.QueryFollow(user, DB)
		c.JSON(200, gin.H{
			"follows": follows,
		})
	})

	r.GET("/follower", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)

		followers := database.QueryFollower(user, DB)
		c.JSON(200, gin.H{
			"follows": followers,
		})
	})

	r.GET("/comment", func(c *gin.Context) {
		setHeader(c)

		gifid := c.DefaultQuery("gifid", "")

		comments := database.QueryComment(gifid, DB)
		c.JSON(200, gin.H{
			"comments": comments,
		})
	})

	r.GET("/user_gifs", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)

		gifs := database.QueryGifs(user, DB)
		c.JSON(200, gin.H{
			"gifs": gifs,
		})
	})

	r.POST("/change_profile", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		Email := c.DefaultPostForm("Email", "")
		FirstName := c.DefaultPostForm("FirstName", "")
		LastName := c.DefaultPostForm("LastName", "")
		Addr := c.DefaultPostForm("Addr", "")
		ZipCode := c.DefaultPostForm("ZipCode", "")
		City := c.DefaultPostForm("City", "")
		Country := c.DefaultPostForm("Country", "")
		About := c.DefaultPostForm("About", "")
		Height := c.DefaultPostForm("Height", "")
		Birthday := c.DefaultPostForm("Birthday", "")

		database.ChangeProfile(user, Email, FirstName, LastName, Addr, ZipCode, City, Country, About, Height, Birthday, DB)
		c.JSON(200, gin.H{
			"status": true,
		})
	})

	// r.GET("/")

	return r
}

func main() {
	cache.OfflineCacheInit()
	r := RouterSet()
	r.Run(":8080")

	// DB := database.ConnectDB()
	// database.Init(DB)
	// database.InsertUser("Admin", "Admin", "", DB)
	// // /Users/saberrrrrrrr/Desktop/spider_info.json
	// gifs := utils.JsonParse("/Users/saberrrrrrrr/Desktop/info_spider.json") //("/Users/saberrrrrrrr/Desktop/backend/info.json")
	// for _, gif := range gifs {
	// 	database.InsertGIF(DB, "Admin", gif.Name, gif.Keyword, "开始的gif", gif.Title)
	// }

	// fmt.Println(cookie.ShaConvert("user0"))

	// goc := cookie.CookieCacheInit()
	// cookie.CookieSet("user0", goc)
	// fmt.Println(cookie.CookieTest(string(cookie.ShaConvert("user0")), goc))
	// res, _:=goc.Get("user0")
	// fmt.Println(res)
}
