package cookie

import(
	gocache "github.com/patrickmn/go-cache"
	"fmt"
	"time"
	"backend/utils"
	"net/http"
)

func CookieCacheInit() *gocache.Cache{
	c:=gocache.New(utils.COOKIE_EXPIRE*time.Second, 600*time.Second)
	return c
}

func CookieSet(user string, c *gocache.Cache){
	c.Set(user, "normal", gocache.DefaultExpiration)
}

func RootCookieSet(user string, c *gocache.Cache){
	c.Set(user, "", gocache.DefaultExpiration)
}

func CookieCheck(req *http.Request, c *gocache.Cache)bool{
	cookie, _:=req.Cookie("user_name")
	if(cookie== nil){
		return false;
	}else{
		_, b:=c.Get(cookie.Value)
		return b;
	}
}

func CookieTest(){
	fmt.Println("test undefined")
}