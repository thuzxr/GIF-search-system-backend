package main

import (
	"backend/cache"
	// "backend/cookie"
	"backend/database"
	"backend/management/login"
	"backend/management/register"
	"backend/management/vericode"
	"backend/ossUpload"
	"backend/recommend"
	"backend/search"
	"backend/upload"
	"backend/utils"
	// "backend/cookie"
	"fmt"

	"time"
	"backend/word"

	"github.com/gin-gonic/gin"
	"github.com/go-ego/gse"
	_ "github.com/go-sql-driver/mysql"
	
    jwt "github.com/dgrijalva/jwt-go"
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

type MyClaims struct {
	User_name string `json:"user_name"`
	Access int `json:"access"`
	jwt.StandardClaims
}

func UserAuth() gin.HandlerFunc{
	return func(c *gin.Context){
		cookie, _:=c.Request.Cookie("token")
		if(cookie!=nil){
			_, status:=ClaimsParse(cookie.Value)
			if(status>=1){
				c.Next();
			}else{
				c.Abort();
				c.JSON(401, gin.H{
					"status": "Unauthorized",
				})
			}
		}else{
			c.Abort();
			c.JSON(401, gin.H{
				"status": "Unauthorized",
			})
		}
	}
}

func UserAntiAuth() gin.HandlerFunc{
	return func(c *gin.Context){
		cookie, _:=c.Request.Cookie("token")
		if(cookie!=nil){
			_, status:=ClaimsParse(cookie.Value)
			if(status>=1){
				c.Abort()
				c.JSON(412, gin.H{
					"status": "Has User Online",
				})
			}else{
				c.Next()
			}
		}else{
			c.Next()
		}
	}
}

func ClaimsParse(tokenString string) (*MyClaims, int){
	var claims *MyClaims
	var status int
	var ok bool
	token, err:= jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.COOKIE_SALT), nil
	})
	fmt.Println("token:",token)
	if err == nil {
		claims, ok = token.Claims.(*MyClaims)
		if ok && token.Valid {
			if(claims.Access == 2) {
				status=2
			}else{
				status=1
			}
		} else {
			fmt.Println("claim not exist", token.Valid)
			status=-1
		}
	} else {
		fmt.Println("err in claim Parse:", err)
		status=0
	}
	return claims, status
}

func TokenSet(c *gin.Context, user string, access int){
	claims:=MyClaims{
		user,
		access,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 3600),
			Issuer: "Gif-Dio",
		},
	}
	fmt.Println(claims)
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err:=token.SignedString([]byte(utils.COOKIE_SALT))
	if err!=nil{
		fmt.Println("err in tokenSet:", err)
		return
	}else{
		c.SetCookie("token", tokenString, 3600, "/", utils.COOKIE_DOMAIN, false, false)
		fmt.Println(tokenString)
	}
}

func TokenTest(user string, access int){
	claims:=MyClaims{
		user,
		access,
		jwt.StandardClaims{
			ExpiresAt:int64(time.Now().Unix() + 3600) ,
			Issuer: "Gif-Dio",
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err:=token.SignedString([]byte(utils.COOKIE_SALT))
	if err!=nil{
		fmt.Println("err in tokenSet:", err)
		return
	}else{
		// c.SetCookie("token", tokenString, 3600, "/", utils.COOKIE_DOMAIN, false, false)
		fmt.Println(tokenString)
	}
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
	names, titles, keywords := search.FastIndexParse()
	// names:=make([]string,0)
	// titles:=make([]string,0)
	// keywords:=make([]string,0)

	fmt.Println(gifs[0])
	var maps map[string]utils.Gifs
	maps = make(map[string]utils.Gifs)
	for _, gif := range gifs {
		maps[gif.Name] = gif
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
	r.POST("/login", //UserAntiAuth(),	
		func(c *gin.Context) {
		setHeader(c)

		// user := c.DefaultQuery("user", "")
		// password := c.DefaultQuery("password", "")
		user:=c.DefaultPostForm("user", "")
		password:=c.DefaultPostForm("password", "")

		fmt.Println(user)
		fmt.Println(password)

		status := login.Login(user, password, DB)
		if(status=="登陆成功！"){
			// c.SetCookie("user_name", string(cookie.ShaConvert(user)), 3600, "/", utils.COOKIE_DOMAIN,  false, false)
			// cookie.CookieSet(user, goc)
			TokenSet(c, user, 1)
			c.JSON(200, gin.H{
				"status": 1,
			})
		}else{
			c.JSON(406, gin.H{
				"status": 0,
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
	
	r.GET("/user_status", func(c *gin.Context){
		res,_:=c.Request.Cookie("token")
		var status int
		var claims *MyClaims
		if res==nil{
			status=0
			claims=&MyClaims{}
		}else{
			tokenString:=res.Value
			fmt.Println(tokenString)
			claims, status=ClaimsParse(tokenString)
		}
		c.JSON(200, gin.H{
			"status": status,
			"claims": claims,
		})
	})

	r.GET("/user_login", func(c *gin.Context){
		user:="user0"
		access:=1
		TokenSet(c, user, access)
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

	// r.Use(UserAuth())

	r.GET("/test", func(c *gin.Context){
		c.JSON(200, gin.H{
			"res":"test",
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

		keyword := c.DefaultPostForm("keyword", "")
		name := c.DefaultPostForm("name", "")
		title := c.DefaultPostForm("title", "")
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

	return r
}

func main() {
	cache.OfflineCacheInit()
	r := RouterSet()
	r.Run(":8080")
	// fmt.Println(cookie.ShaConvert("user0"))

	// goc := cookie.CookieCacheInit()
	// cookie.CookieSet("user0", goc)
	// fmt.Println(cookie.CookieTest(string(cookie.ShaConvert("user0")), goc))
	// res, _:=goc.Get("user0")
	// fmt.Println(res)
}
