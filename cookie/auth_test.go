package cookie

import(	
	"testing"
	"regexp"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
)

func getToken(w *httptest.ResponseRecorder) string{
	exp1:=regexp.MustCompile(`token=(.*?);`)
	cookie_header:=w.HeaderMap["Set-Cookie"][0]
	res:=exp1.FindAllStringSubmatch(cookie_header,-1)
	if(len(res)<1){
		return "";
	}else if(len(res[0])<1){
		return "";
	}else{
		return string(res[0][1])
	}	
}

func TestAuth(t *testing.T){
	r:=gin.Default()

	r.GET("/test", func(c *gin.Context){
		TokenSet(c, "Dio", 1)
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/auth", UserAuth(), func(c *gin.Context){
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/anti", UserAntiAuth(), func(c* gin.Context){
		c.JSON(200, gin.H{
			"status": "succeed",
		})
	})

	r.GET("/test_cookie",func(c *gin.Context){
		c.String(200, Getusername(c))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	w=httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/anti", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w=httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)
	jwtoken:=getToken(w);
	assert.NotEqual(t, jwtoken, "")

	w=httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/test_cookie", nil)
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:jwtoken,
	})
	r.ServeHTTP(w, req)
	assert.Equal(t, string(w.Body.Bytes()), "Dio")	

	claims,status:=ClaimsParse("")
	assert.Equal(t, status, 0)
	claims,status=ClaimsParse(jwtoken)
	assert.Equal(t, claims.Access, status)

	w=httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/auth", nil)
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:jwtoken,
	})
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w=httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/anti", nil)
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:jwtoken,
	})
	r.ServeHTTP(w, req)
	assert.Equal(t, 412, w.Code)
}
