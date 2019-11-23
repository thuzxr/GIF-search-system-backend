package main

import (
	"backend/cache"
	"sort"
	"strings"

	"backend/cookie"
	"backend/database"
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
	"backend/router"
	goini "github.com/clod-moon/goconf"
)

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

	// AdSearch_Enabled := word.DataCheck()
	// _ := AdSearch_Enabled
	AdSearch_Enabled := false

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
						if ! b {
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
		router.SetHeader(c)

		msg := c.DefaultQuery("msg", "000")
		fmt.Println(msg)
		c.JSON(200, gin.H{
			"message": "hello world! --sent by GO",
		})
	})
	r.GET("/search", func(c *gin.Context) {
		time0 := time.Now()
		router.SetHeader(c)

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

		if len(match) == 0 {
			c.JSON(200, gin.H{
				utils.STATUS: "failed",
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
					utils.STATUS:   utils.SUCCEED,
					utils.RESULT:   result_match,
					"like_num": result_score,
					"time":     t0,
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
					utils.STATUS:   utils.SUCCEED,
					utils.RESULT:   match,
					"like_num": likes_num,
					"time":     t0,
				})
			}
		}
	})
	
	router.ManageRouterSet(r, DB, likes_u2g)

	router.CaptchaRouterSet(r)

	//Routers with Auth

	r.Use(cookie.UserAuth())

	router.OtherRouterSet(r, DB)

	r.POST("/verify", func(c *gin.Context) {
		router.SetHeader(c)
		veriName := c.DefaultPostForm("name", "")
		veriNames := strings.Split(veriName, " ")
		var veriName0 string
		for i := range veriNames {
			if len(veriNames[i]) == 0 {
				continue
			}
			veriName0 = veriNames[i]
			database.VerifyGIF(DB, veriNames[i])
			fmt.Println("verifing gif ", veriName0, "@@")
			ossUpload.OssMove(veriName0)
			fmt.Println("verifing gif ", veriName0, "@@")
		}
		c.JSON(200, gin.H{
			utils.STATUS: utils.SUCCEED,
		})
		ch_gifUpdate <- true
		fmt.Println("verifing over")
	})

	router.VerifyRouterSet(r, DB)

	r.POST("/remove", func(c *gin.Context) {
		router.SetHeader(c)
		name := c.DefaultPostForm("name", "")
		database.DeleteGif(name, DB)
		c.JSON(200, gin.H{
			utils.STATUS: utils.SUCCEED,
		})
	})

	r.GET("/recommend", func(c *gin.Context) {
		router.SetHeader(c)

		user := cookie.Getusername(c)

		recom, ok := userCF[user]
		//fmt.Println(recom)
		if !ok || len(recom) == 0 {
			c.JSON(200, gin.H{
				utils.STATUS: utils.SUCCEED,
				utils.RESULT: gifs[:10],
			})
		} else {
			var rets []utils.Gifs
			for _, rec := range recom {
				rets = append(rets, gifs[maps[rec]])
			}
			c.JSON(200, gin.H{
				utils.STATUS: utils.SUCCEED,
				utils.RESULT: rets,
			})
		}

		//name := c.DefaultQuery("name", "")
		//recommend_gifs := recommend.Recommend(gifs[maps[name]], gifs)
		//c.JSON(200, gin.H{
		//	utils.STATUS: utils.SUCCEED,
		//	utils.RESULT: recommend_gifs,
		//})
	})


	router.ProfileRouterSet(r, DB)

	r.GET("/favor", func(c *gin.Context) {
		router.SetHeader(c)

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
			utils.RESULT: results,
		})
	})

	router.FavorRouterSet(r, likes, likes_u2g)

	// r.GET("/follow", func(c *gin.Context) {
	// 	router.SetHeader(c)

	// 	// user := c.DefaultQuery("user", "")
	// 	user := cookie.Getusername(c)

	// 	follows := database.QueryFollow(user, DB)
	// 	c.JSON(200, gin.H{
	// 		"follows": follows,
	// 	})
	// })

	// r.GET("/follower", func(c *gin.Context) {
	// 	router.SetHeader(c)

	// 	// user := c.DefaultQuery("user", "")
	// 	user := cookie.Getusername(c)

	// 	followers := database.QueryFollower(user, DB)
	// 	c.JSON(200, gin.H{
	// 		"follows": followers,
	// 	})
	// })

	// r.GET("/comment", func(c *gin.Context) {
	// 	router.SetHeader(c)

	// 	gifid := c.DefaultQuery("gifid", "")

	// 	comments := database.QueryComment(gifid, DB)
	// 	c.JSON(200, gin.H{
	// 		"comments": comments,
	// 	})
	// })

	// r.GET("/user_gifs", func(c *gin.Context) {
	// 	router.SetHeader(c)

	// 	// user := c.DefaultQuery("user", "")
	// 	user := cookie.Getusername(c)

	// 	gifs := database.QueryGifs(user, DB)
	// 	c.JSON(200, gin.H{
	// 		"gifs": gifs,
	// 	})
	// })

	// r.POST("/like", func(c *gin.Context) {
	// 	router.SetHeader(c)
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

	return r
}

func LoadTls() gin.HandlerFunc {
	
	conf:=goini.InitConfig("settings.ini")
	hostAddr:=conf.GetValue("ssl","sslhost")
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     hostAddr,
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
	// // r.Run(":8080")
	r.Use(LoadTls())
	r.RunTLS(":8080", "/etc/nginx/1_www.gifxiv.com_bundle.crt", "/etc/nginx/2_www.gifxiv.com.key")

	// DB := database.ConnectDB()
	// database.Init(DB)
	// database.InsertUser("Admin", "Admin", DB)
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
