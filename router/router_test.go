package router

import(
	"testing"
	
	"github.com/stretchr/testify/assert"
	"net/http"
	// "net/url"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"backend/database"
	"fmt"
	// "strings"
)

func TestRouters(t *testing.T){
	r:=gin.Default()
	DB := database.ConnectDB()
	
	ProfileRouterSet(r, DB)
	w:=httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	r.ServeHTTP(w, req)
	w=httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/change_profile", nil)
	r.ServeHTTP(w, req)
	fmt.Println(string(w.Body.Bytes()))
	assert.Equal(t, 200, w.Code)

	// test_name:="4fd32a6fae93404a956129260ec0a606"
	_, _, _, likes, likes_u2g := database.LoadAll(DB)
	FavorRouterSet(r, likes, likes_u2g)
	req, _ = http.NewRequest(http.MethodPost, "/insert_favor", nil)
	r.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodPost, "/delete_favor", nil)
	r.ServeHTTP(w, req)

	VerifyRouterSet(r, DB)
	req, _ = http.NewRequest(http.MethodPost, "/remove_verify", nil)
	r.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodGet, "/toBeVerify", nil)
	r.ServeHTTP(w, req)

	OtherRouterSet(r, DB)
	req, _ = http.NewRequest(http.MethodGet, "/logout", nil)
	r.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodPost, "/upload", nil)
	r.ServeHTTP(w, req)

	CaptchaRouterSet(r)
	req, _ = http.NewRequest(http.MethodGet, "/refresh_veri", nil)
	r.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodGet, "/get_veri/:captchId", nil)
	r.ServeHTTP(w, req)

	ManageRouterSet(r, DB, likes_u2g)
	req, _ = http.NewRequest(http.MethodPost, "/login", nil)
	r.ServeHTTP(w, req)
	// data:=url.Values{}
	// data.Set("user", "jojo")
	// data.Set("password", "jojo")
	// req, _ = http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	// r.ServeHTTP(w, req)

	req, _ = http.NewRequest(http.MethodPost, "/register", nil)
	r.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodGet, "/user_status", nil)
	r.ServeHTTP(w, req)
}
