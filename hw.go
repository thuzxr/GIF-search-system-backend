package main

import (
	"backend/cache"
	"sort"
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
	"strconv"
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
	_, _, gifs, likes, likes_u2g := database.LoadAll(DB)
	userCF := recommend.UserCF(likes, likes_u2g)
	gif_proto := utils.JsonParse("info_old_recommend.json")
	fmt.Println("gif proto", gif_proto[1])

	AdSearch_Enabled := word.DataCheck()
	// _ := AdSearch_Enabled
	// AdSearch_Enabled := false
	// AdSearch_Activated := AdSearch_Enabled

	var gif2vec map[string][][]uint8
	var word2vec map[string][]uint8
	var re_idx []string
	var vec_h [][]uint64
	var seg gse.Segmenter

	fmt.Println("OssUpdating")
	ossUpload.OssUpdate(gifs)
	fmt.Println("OssUpdated")
	if len(gifs) > 0 {
		fmt.Println(gifs[0].Oss_url)
	} else {
		fmt.Println("###### WARNING #######       GIF LOAD FAILED , Checkout your database")
	}

	var maps map[string]int
	maps = make(map[string]int)
	for i := range gifs {
		maps[gifs[i].Name] = i
	}

	var rec_tmp []int
	for i := range gif_proto {
		rec_tmp = make([]int, 0)
		for j := range gif_proto[i].Recommend {
			rec_tmp = append(rec_tmp, maps[gif_proto[gif_proto[i].Recommend[j]].Name])
		}
		gifs[maps[gif_proto[i].Name]].Recommend = rec_tmp
	}

	for i := range gifs {
		rec_tmp = make([]int, 0)
		if len(gifs[i].Recommend) == 0 {
			for j := 0; j < 10; j++ {
				rec_tmp = append(rec_tmp, ((int)(rand.Int31()) % (len(gifs))))
			}
			gifs[i].Recommend = rec_tmp
		}
	}
	// fmt.Println("## gif 0:",gifs[1], gifs[0])

	if AdSearch_Enabled {
		fmt.Println("Advanced Searching Enabled")
		word2vec = word.WordToVecInit()
		re_idx, gif2vec, vec_h = word.RankSearchInit()
		seg.LoadDict()
	} else {
		fmt.Println("Index not found, Advanced Searching Disabled")
	}

	fmt.Println("total gifs size ", len(gifs))

	ch_gifUpdate := make(chan bool)
	go func() {
		for {
			select {
			case <-ch_gifUpdate:
				_, infos2, gifs2, likes2, likes_u2g2 := database.LoadAll(DB)
				maps2 := make(map[string]int)
				for i := range gifs2 {
					res, b := maps[gifs2[i].Name]
					if b {
						gifs2[i].Oss_url = gifs[res].Oss_url
					} else {
						gifs2[i].Oss_url = ossUpload.OssSignLink(gifs2[i], 3600)
					}
					maps2[gifs2[i].Name] = i
				}
				if AdSearch_Enabled {
					// veci:=word.WortToVec(gifs)
					for i := range gifs2 {
						_, b := gif2vec[gifs2[i].Name]
						if b == false {
							veci, vechi, re_idxi := word.GifToVec(gifs2[i], seg, word2vec)
							gif2vec[gifs2[i].Name] = veci
							vec_h = append(vec_h, vechi...)
							re_idx = append(re_idx, re_idxi...)
							fmt.Println("discovered new gif in ad search ", gifs2[i].Name)
						}
					}
				}
				// users = users2
				gifs = gifs2
				_ = infos2
				likes = likes2
				maps = maps2
				likes_u2g = likes_u2g2
				fmt.Println("gif updates here")
				fmt.Println("total gifs size ", len(gifs))
				fmt.Println("spec here ", maps["0000aaaa"])

				ch_gifUpdate <- false
				break
			default:
				break
			}
		}
	}()

	// fmt.Println(gifs[0])

	go func() {
		for {
			time.Sleep(50 * time.Minute)
			// time.Sleep(30*time.Second)
			fmt.Println("Oss&likes Updating")
			ossUpload.OssUpdate(gifs)
			database.UpdateLikes(likes, DB)
			fmt.Println(gifs[0].Oss_url)
			fmt.Println("Oss&likes Updated")
		}
	}()

	go func() {
		for {
			time.Sleep(time.Minute)
			fmt.Println("user_CF Updating")
			userCF = recommend.UserCF(likes, likes_u2g)
			fmt.Println("user_CF Updated")
		}
	}()

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
		time0 := time.Now()
		setHeader(c)

		// time0:=time.Now()
		keyword := c.DefaultQuery("key", "UNK")
		typ := c.DefaultQuery("type", "L")
		rank_type := c.DefaultQuery("rank_type", "Sim")
		edg := c.DefaultQuery("edge", "200")

		edg0, err := strconv.ParseInt(edg, 10, 64)
		if err != nil {
			edg0 = 10
		} else if edg0 < 1 || edg0 > 10 {
			edg0 = 10
		}
		edg0 = edg0*10 + 150
		HAM_EDGE := uint64(edg0)

		keyw0 := typ + edg + keyword
		res, finded := m[keyw0]
		var match []utils.Gifs
		// fmt.Println(time.Since(time0))
		if finded {
			match = make([]utils.Gifs, len(res))
			for i := range res {
				match[i] = gifs[maps[res[i].Name]]
			}
			m[keyw0] = match
			fmt.Println("Hit Cache " + keyword)
		} else {
			if typ == "H" {
				res := word.RankSearch(keyword, word2vec, gif2vec, vec_h, re_idx, seg, HAM_EDGE)
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
			m[keyw0] = match
			go cache.OfflineCacheAppend(keyword, match, typ, edg)
		}
		t0 := int64(time.Since(time0) / time.Nanosecond)
		if t0 < 10 {
			t0 = 1
		}
		// for i := 0; i < len(match); i++ {
		// 	match[i].Oss_url = ossUpload.OssSignLink(match[i], 3600)
		// }
		// fmt.Println(time.Since(time0))

		if len(match) == 0 {
			c.JSON(200, gin.H{
				"status": "failed",
				"time":   t0,
			})
		} else {
			if rank_type == "Sim" {
				var match_likes []utils.Like_based_sort
				for i := range match {
					_, ok := likes[match[i].Name]
					var match_like utils.Like_based_sort
					if ok {
						match_like = utils.Like_based_sort{Gif: match[i], Like: len(likes[match[i].Name])}
					} else {
						match_like = utils.Like_based_sort{Gif: match[i], Like: 0}
					}
					match_likes = append(match_likes, match_like)
				}

				sort.Sort(utils.LikeSlice(match_likes))

				var result_match []utils.Gifs
				var result_score []int
				for i := range match_likes {
					result_match = append(result_match, match_likes[i].Gif)
					result_score = append(result_score, match_likes[i].Like)
				}

				c.JSON(200, gin.H{
					"status":   "succeed",
					"result":   result_match,
					"like_num": result_score,
					"time": t0,
				})
			} else {
				var likes_num []int
				for i := range match {
					_, ok := likes[match[i].Name]
					if ok {
						likes_num = append(likes_num, len(likes[match[i].Name]))
					} else {
						likes_num = append(likes_num, 0)
					}
				}

				c.JSON(200, gin.H{
					"status":   "succeed",
					"result":   match,
					"like_num": likes_num,
					"time":t0,
				})
			}
		}
	})

	r.POST("/login", cookie.UserAntiAuth(), func(c *gin.Context) {
		setHeader(c)

		user := c.DefaultPostForm("user", "")
		password := c.DefaultPostForm("password", "")

		status := login.Login(user, password, DB)

		if status != -1 {
			cookie.TokenSet(c, user, status)
			// favors := database.QueryFavor(user, DB)
			favors := likes_u2g[user]
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
		for i := range res {
			res[i].OSSURL = ossUpload.OssSignLink_Verify(utils.Gifs{
				Name: res[i].GifId,
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
		veriNames := strings.Split(veriName, " ")
		var veriName0 string
		for i := range veriNames {
			if(len(veriNames[i])==0){
				continue
			}
			veriName0=veriNames[i]
			database.VerifyGIF(DB, veriNames[i])
			fmt.Println("verifing gif ", veriName0, "@@")
			ossUpload.OssMove(veriName0)
			fmt.Println("verifing gif ", veriName0, "@@")
		}
		c.JSON(200, gin.H{
			"status": "succeed",
		})
		ch_gifUpdate <- true
		fmt.Println("verifing over")
	})

	r.POST("/remove_verify", func(c *gin.Context) {
		setHeader(c)
		name := c.DefaultPostForm("name", "")
		removeNames := strings.Split(name, " ")
		for i := range removeNames {
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

		user := cookie.Getusername(c)

		recom, ok := userCF[user]
		//fmt.Println(recom)
		if !ok || len(recom) == 0 {
			c.JSON(200, gin.H{
				"status": "succeed",
				"result": gifs[:10],
			})
		} else {
			var rets []utils.Gifs
			for _, rec := range recom {
				rets = append(rets, gifs[maps[rec]])
			}
			c.JSON(200, gin.H{
				"status": "succeed",
				"result": rets,
			})
		}

		 //name := c.DefaultQuery("name", "")
		 //recommend_gifs := recommend.Recommend(gifs[maps[name]], gifs)
		 //c.JSON(200, gin.H{
		//	"status": "succeed",
		//	"result": recommend_gifs,
		 //})
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

		// favors := database.QueryFavor(user, DB)
		favors := likes_u2g[user]
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
		// favors := database.InsertFavor(user, gifid, DB)
		likes[gifid] = append(likes[gifid], user)
		likes_u2g[user] = append(likes_u2g[user], gifid)

		c.JSON(200, gin.H{
			"status": "收藏成功",
		})
	})

	r.POST("/delete_favor", func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		gifid_string := c.DefaultPostForm("GifId", "")
		gifids := strings.Split(gifid_string, " ")
		// favors := database.DeleteFavor(user, gifids, DB)
		for _, gifid := range gifids {
			for j, usr := range likes[gifid] {
				if usr == user {
					likes[gifid] = append(likes[gifid][:j], likes[gifid][j+1:]...)
					break
				}
			}
		}

		for _, gifid := range gifids {
			for i, checkgif := range likes_u2g[user] {
				if gifid == checkgif {
					likes_u2g[user] = append(likes_u2g[user][:i], likes_u2g[user][i+1:]...)
					break
				}
			}
		}

		c.JSON(200, gin.H{
			"status": "删除成功",
		})
	})

	// r.GET("/follow", func(c *gin.Context) {
	// 	setHeader(c)

	// 	// user := c.DefaultQuery("user", "")
	// 	user := cookie.Getusername(c)

	// 	follows := database.QueryFollow(user, DB)
	// 	c.JSON(200, gin.H{
	// 		"follows": follows,
	// 	})
	// })

	// r.GET("/follower", func(c *gin.Context) {
	// 	setHeader(c)

	// 	// user := c.DefaultQuery("user", "")
	// 	user := cookie.Getusername(c)

	// 	followers := database.QueryFollower(user, DB)
	// 	c.JSON(200, gin.H{
	// 		"follows": followers,
	// 	})
	// })

	// r.GET("/comment", func(c *gin.Context) {
	// 	setHeader(c)

	// 	gifid := c.DefaultQuery("gifid", "")

	// 	comments := database.QueryComment(gifid, DB)
	// 	c.JSON(200, gin.H{
	// 		"comments": comments,
	// 	})
	// })

	// r.GET("/user_gifs", func(c *gin.Context) {
	// 	setHeader(c)

	// 	// user := c.DefaultQuery("user", "")
	// 	user := cookie.Getusername(c)

	// 	gifs := database.QueryGifs(user, DB)
	// 	c.JSON(200, gin.H{
	// 		"gifs": gifs,
	// 	})
	// })

	// r.POST("/like", func(c *gin.Context) {
	// 	setHeader(c)
	// 	user := cookie.Getusername(c)
	// 	GifId := c.DefaultPostForm("GifId", "")
	// 	contains := false
	// 	_, ok := likes[GifId]
	// 	if ok {
	// 		for i := range likes[GifId] {
	// 			if likes[GifId][i] == user {
	// 				contains = true
	// 				likes[GifId] = append(likes[GifId][:i], likes[GifId][i+1:]...)
	// 				break
	// 			}
	// 		}
	// 		if !contains {
	// 			likes[GifId] = append(likes[GifId], user)
	// 		}
	// 	} else {
	// 		likes[GifId] = make([]string, 0)
	// 		likes[GifId] = append(likes[GifId], user)
	// 	}
	// })

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
	// r.Run(":8080")
	r.Use(LoadTls())
	r.RunTLS(":8080", "/etc/nginx/1_www.gifxiv.com_bundle.crt", "/etc/nginx/2_www.gifxiv.com.key")

	// DB := database.ConnectDB()
	// database.Init(DB)
	// database.InsertUser("Admin", "Admin", "", DB)
	// // // /Users/saberrrrrrrr/Desktop/spider_info.json
	// gifs := utils.JsonParse("/Users/saberrrrrrrr/Desktop/backend/info_old_recommend.json")
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
