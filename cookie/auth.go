package cookie

import (
	"backend/utils"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MyClaims struct {
	User_name string `json:"user_name"`
	Access    int    `json:"access"`
	jwt.StandardClaims
}

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("token")
		if cookie != nil {
			_, status := ClaimsParse(cookie.Value)
			if status >= 1 {
				c.Next()
			} else {
				c.Abort()
				c.JSON(401, gin.H{
					"status": "Unauthorized",
				})
			}
		} else {
			c.Abort()
			c.JSON(401, gin.H{
				"status": "Unauthorized",
			})
		}
	}
}

func UserAntiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("token")
		if cookie != nil {
			_, status := ClaimsParse(cookie.Value)
			if status >= 1 {
				c.Abort()
				c.JSON(412, gin.H{
					"status": "Has User Online",
				})
			} else {
				c.Next()
			}
		} else {
			c.Next()
		}
	}
}

func ClaimsParse(tokenString string) (*MyClaims, int) {
	var claims *MyClaims
	var status int
	var ok bool
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.COOKIE_SALT), nil
	})
	fmt.Println("token:", token)
	if err == nil {
		claims, ok = token.Claims.(*MyClaims)
		if ok && token.Valid {
			if claims.Access == 2 {
				status = 2
			} else {
				status = 1
			}
		} else {
			fmt.Println("claim not exist", token.Valid)
			status = -1
		}
	} else {
		fmt.Println("err in claim Parse:", err)
		status = 0
	}
	return claims, status
}

func TokenSet(c *gin.Context, user string, access int) {
	claims := MyClaims{
		user,
		access,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 3600),
			Issuer:    "Gif-Dio",
		},
	}
	fmt.Println(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(utils.COOKIE_SALT))
	if err != nil {
		fmt.Println("err in tokenSet:", err)
		return
	} else {
		c.SetCookie("token", tokenString, 3600, "/", utils.COOKIE_DOMAIN, false, false)
		fmt.Println(tokenString)
	}
}

func TokenTest(user string, access int) {
	claims := MyClaims{
		user,
		access,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 3600),
			Issuer:    "Gif-Dio",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(utils.COOKIE_SALT))
	if err != nil {
		fmt.Println("err in tokenSet:", err)
		return
	} else {
		// c.SetCookie("token", tokenString, 3600, "/", utils.COOKIE_DOMAIN, false, false)
		fmt.Println(tokenString)
	}
}

func Getusername(c *gin.Context) string {
	cookie, _ := c.Request.Cookie("token")
	if cookie != nil {
		tokenString := cookie.Value
		var claims *MyClaims
		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.COOKIE_SALT), nil
		})
		if err == nil {
			claims, _ = token.Claims.(*MyClaims)
			return claims.User_name
		} else {
			return ""
		}
	}
	return ""
}

func RootAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("token")
		if cookie != nil {
			_, status := ClaimsParse(cookie.Value)
			if status >= 2 {
				c.Next()
			} else {
				c.Abort()
				c.JSON(401, gin.H{
					"status": "Not Root",
				})
			}
		} else {
			c.Abort()
			c.JSON(401, gin.H{
				"status": "Unauthorized",
			})
		}
	}
}
