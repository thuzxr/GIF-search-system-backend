package router

import (
	"github.com/gin-gonic/gin"

	"backend/cookie"
	"backend/database"
	"backend/management/login"
	"backend/management/register"
	"backend/management/vericode"
	"backend/ossUpload"
	"backend/utils"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func ProfileRouterSet(r *gin.Engine, DB *sql.DB) {
	r.GET("/profile", func(c *gin.Context) {
		SetHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		profile := database.QueryProfile(user, DB)
		c.JSON(200, gin.H{
			"Email":         profile[0],
			utils.FIRSTNAME: profile[1],
			utils.LASTNAME:  profile[2],
			"Addr":          profile[3],
			utils.ZIPCODE:   profile[4],
			"City":          profile[5],
			utils.COUNTRY:   profile[6],
			"About":         profile[7],
			utils.HEIGHT:    profile[8],
			utils.BIRTHDAY:  profile[9],
		})
	})

	r.POST("/change_profile", func(c *gin.Context) {
		SetHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		Email := c.DefaultPostForm("Email", "")
		FirstName := c.DefaultPostForm(utils.FIRSTNAME, "")
		LastName := c.DefaultPostForm(utils.LASTNAME, "")
		Addr := c.DefaultPostForm("Addr", "")
		ZipCode := c.DefaultPostForm(utils.ZIPCODE, "")
		City := c.DefaultPostForm("City", "")
		Country := c.DefaultPostForm(utils.COUNTRY, "")
		About := c.DefaultPostForm("About", "")
		Height := c.DefaultPostForm(utils.HEIGHT, "")
		Birthday := c.DefaultPostForm(utils.BIRTHDAY, "")

		profile := utils.Profile{Email: Email, FirstName: FirstName, LastName: LastName, Addr: Addr, ZipCode: ZipCode, City: City,
			Country: Country, About: About, Height: Height, Birthday: Birthday}

		database.ChangeProfile(user, profile, DB)
		c.JSON(200, gin.H{
			utils.STATUS: true,
		})
	})
}

func SetHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
}

func FavorRouterSet(r *gin.Engine, likes, likes_u2g map[string][]string) {
	r.POST("/insert_favor", func(c *gin.Context) {
		SetHeader(c)

		// user := c.DefaultQuery("user", "")
		user := cookie.Getusername(c)
		gifid := c.DefaultPostForm("GifId", "")
		// favors := database.InsertFavor(user, gifid, DB)
		likes[gifid] = append(likes[gifid], user)
		likes_u2g[user] = append(likes_u2g[user], gifid)

		c.JSON(200, gin.H{
			utils.STATUS: "收藏成功",
		})
	})

	r.POST("/delete_favor", func(c *gin.Context) {
		SetHeader(c)

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
			utils.STATUS: "删除成功",
		})
	})
}

func VerifyRouterSet(r *gin.Engine, DB *sql.DB) {

	r.POST("/remove_verify", func(c *gin.Context) {
		SetHeader(c)
		name := c.DefaultPostForm("name", "")
		removeNames := strings.Split(name, " ")
		for i := range removeNames {
			database.RemoveVerify(DB, removeNames[i])
		}
		c.JSON(200, gin.H{
			utils.STATUS: utils.SUCCEED,
		})
	})
	r.GET("/toBeVerify", func(c *gin.Context) {
		SetHeader(c)

		res := database.GetToVerifyGIF(DB)
		for i := range res {
			res[i].OSSURL = ossUpload.OssSignLink_Verify(utils.Gifs{
				Name: res[i].GifId,
			}, 3600)
		}
		c.JSON(200, gin.H{
			utils.STATUS: utils.SUCCEED,
			utils.RESULT: res,
		})
	})
}

func OtherRouterSet(r *gin.Engine, DB *sql.DB) {
	r.GET("/logout", func(c *gin.Context) {
		SetHeader(c)
		c.SetCookie("token", "", -1, "/", utils.COOKIE_DOMAIN, false, false)
		c.JSON(200, gin.H{
			utils.STATUS: "Logout",
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		SetHeader(c)

		user := cookie.Getusername(c)
		info := c.DefaultPostForm("info", "")
		keyword := c.DefaultPostForm("keyword", "")
		name := c.DefaultPostForm("name", "")
		title := c.DefaultPostForm("title", "")
		database.InsertUnderVerifyGIF(DB, user, name, keyword, info, title)
		c.JSON(200, gin.H{
			utils.STATUS: utils.SUCCEED,
		})
	})
}

func CaptchaRouterSet(r *gin.Engine) {

	r.GET("/refresh_veri", func(c *gin.Context) {
		SetHeader(c)
		vericode.Getvericode(c)
	})

	r.GET("/get_veri/:captchId", func(c *gin.Context) {
		SetHeader(c)
		vericode.Genvericode(c)
	})
}

func ManageRouterSet(r *gin.Engine, DB *sql.DB, likes_u2g map[string][]string) {
	r.POST("/login", cookie.UserAntiAuth(), func(c *gin.Context) {
		SetHeader(c)

		user := c.DefaultPostForm("user", "")
		password := c.DefaultPostForm("password", "")

		status := login.Login(user, password, DB)

		if status != -1 {
			cookie.TokenSet(c, user, status)
			// favors := database.QueryFavor(user, DB)
			favors := likes_u2g[user]
			profile := database.QueryProfile(user, DB)

			c.JSON(200, gin.H{
				utils.STATUS:    status,
				"Email":         profile[0],
				utils.FIRSTNAME: profile[1],
				utils.LASTNAME:  profile[2],
				"Addr":          profile[3],
				utils.ZIPCODE:   profile[4],
				"City":          profile[5],
				utils.COUNTRY:   profile[6],
				"About":         profile[7],
				utils.HEIGHT:    profile[8],
				utils.BIRTHDAY:  profile[9],
				"favor":         favors,
			})
		} else {
			c.JSON(406, gin.H{
				utils.STATUS: -1,
			})
		}
	})

	r.POST("/register", func(c *gin.Context) {
		SetHeader(c)
		status := register.Register(c, DB)
		c.JSON(200, gin.H{
			utils.STATUS: status,
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
			claims, status = cookie.ClaimsParse(tokenString)
		}
		c.JSON(200, gin.H{
			utils.STATUS: status,
			"claims":     claims,
		})
	})
}
