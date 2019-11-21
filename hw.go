package main

import (
	"backend/cache"
	"strings"

	"backend/cookie"
	"backend/database"
	"backend/management/login"
	"backend/management/register"
	"backend/management/vericode"
	"backend/ossUpload"
	"backend/recommend"
	"backend/search"

	"backend/utils"
	"math/rand"

	"fmt"
	"time"

	"backend/word"

	"github.com/gin-gonic/gin"
	"github.com/go-ego/gse"
	_ "github.com/go-sql-driver/mysql"
	"github.com/unrolled/secure"
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

	// gifs := utils.JsonParse("info.json")
	users, _, gifs := database.LoadAll(DB)
	gif_proto:=utils.JsonParse("info_old_recommend.json")
	fmt.Println("gif proto", gif_proto[1])

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

	var maps map[string]int
	maps = make(map[string]int)
	for i := range gifs {
		maps[gifs[i].Name] = i
	}

	var rec_tmp []int
	for i:=range(gif_proto){
		rec_tmp=make([]int, 0)
		for j:=range(gif_proto[i].Recommend){
			rec_tmp=append(rec_tmp, maps[gif_proto[gif_proto[i].Recommend[j]].Name])
		}
		gifs[maps[gif_proto[i].Name]].Recommend=rec_tmp
	}

	for i:=range(gifs){
		rec_tmp=make([]int, 0)
		if(len(gifs[i].Recommend)==0){
			for j:=0;j<10;j++{
				rec_tmp=append(rec_tmp, ((int)(rand.Int31())%(len(gifs))))
			}
			gifs[i].Recommend=rec_tmp
		}
	}
	fmt.Println("## gif 0:",gifs[1], gifs[0])

	go func() {
		for {
			time.Sleep(50 * time.Minute)
			// time.Sleep(30*time.Second)
			fmt.Println("OssUpdating")
			ossUpload.OssUpdate(gifs)
			fmt.Println(gifs[0].Oss_url)
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

	fmt.Println("total gifs size ", len(users))

	ch_gifUpdate := make(chan bool)
	go func() {
		for {
			select {
			case <-ch_gifUpdate:
				users2, infos2, gifs2 := database.LoadAll(DB)
				maps2 := make(map[string]int)
				for i := range gifs2 {
					res,b:=maps[gifs2[i].Name]
					if(b){
						gifs2[i].Oss_url=gifs[res].Oss_url
					}else{
						gifs2[i].Oss_url=ossUpload.OssSignLink(gifs2[i], 3600);
					}
					maps2[gifs2[i].Name] = i
				}
				if AdSearch_Enabled {
					// veci:=word.WortToVec(gifs)
					for i:=range(gifs2){
						_,b:=gif2vec[gifs2[i].Name]
						if(b==false){
							veci, vechi, re_idxi:=word.GifToVec(gifs2[i], seg, word2vec)
							gif2vec[gifs2[i].Name]=veci
							vec_h=append(vec_h, vechi...)
							re_idx=append(re_idx,re_idxi...)
						}
					}
				}
				users = users2
				gifs = gifs2
				_ = infos2
				maps=maps2
				fmt.Println("gif updates here")
				fmt.Println("total gifs size ", len(gifs))

				ch_gifUpdate <- false
				break
			default:
				break
			}
		}
	}()

	fmt.Println(gifs[0])

	m := cache.OfflineCacheReload()

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
			match = make([]utils.Gifs, len(res))
			for i := range res {
				match[i] = gifs[maps[res[i].Name]]
			}
			m[keyword] = match
			fmt.Println("Hit Cache " + keyword)
		} else {
			if AdSearch_Enabled {
				res := word.RankSearch(keyword, word2vec, gif2vec, vec_h, re_idx, seg)
				match = make([]utils.Gifs, len(res))
				for i := range res {
					match[i] = gifs[maps[res[i]]]
				}
			} else {
				match0 := search.SimpleSearch(keyword, gifs)
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
		database.InsertUnderVerifyGIF(DB, user, name, keyword, info, title)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/toBeVerify", func(c *gin.Context) {
		setHeader(c)

		res := database.GetToVerifyGIF(DB)
		for i:=range(res){
			res[i].OSSURL=ossUpload.OssSignLink_Verify(utils.Gifs{
				Name:res[i].GifId,
			}, 3600)
		}
		c.JSON(200, gin.H{
			"status": "succeed",
			"result": res,
		})
	})

	r.POST("/verify", func(c *gin.Context) {
		setHeader(c)
		veriName := c.DefaultPostForm("name", "")
		veriNames:=strings.Split(veriName, " ")
		for i:=range(veriNames){
			database.VerifyGIF(DB, veriNames[i])
			ossUpload.OssMove(veriNames[i])
		}
		ch_gifUpdate<-true;
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.POST("/remove_verify", func(c *gin.Context) {
		setHeader(c)
		name := c.DefaultPostForm("name", "")
		removeNames:=strings.Split(name, " ")
		for i:=range(removeNames){
			database.RemoveVerify(DB, removeNames[i])
		}
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

		var results []utils.Gifs
		for favor_id := range favors {
			favor := favors[favor_id]
			results = append(results, gifs[maps[favor]])
		}
		c.JSON(200, gin.H{
			"result": results,
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

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "49.233.71.202:8080",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			//如果出现错误，请不要继续。
			fmt.Println(err)
			return
		}
		// 继续往下处理
		c.Next()
	}
}

func main() {
	cache.OfflineCacheInit()
	r := RouterSet()
	r.Run()
	// r.Use(LoadTls())
	// r.RunTLS(":8080", "/etc/nginx/1_www.gifxiv.com_bundle.crt", "/etc/nginx/2_www.gifxiv.com.key")

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
