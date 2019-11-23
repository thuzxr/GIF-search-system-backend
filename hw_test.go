package main

import (
	"backend/cache"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/utils"
	"backend/cookie"
	"time"
	jwt "github.com/dgrijalva/jwt-go"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRouter(t *testing.T) {
	router := RouterSet()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	message := jsoniter.Get(w.Body.Bytes(), "message").ToString()
	assert.Equal(t, "hello world! --sent by GO", message)
}

func TestSearchRouter(t *testing.T) {
	cache.OfflineCacheInit()
	cache.OfflineCacheClear()
	router := RouterSet()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/search", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	status := jsoniter.Get(w.Body.Bytes(), "status").ToString()
	assert.Equal(t, status, "failed")

	req, _ = http.NewRequest(http.MethodGet, "/search?key=吐出来", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	status = jsoniter.Get(w.Body.Bytes(), "status").ToString()
	assert.Equal(t, status, "failed")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	status = jsoniter.Get(w.Body.Bytes(), "status").ToString()
	assert.Equal(t, status, "failed")
	req, _ = http.NewRequest(http.MethodGet, "/search?key=吐出来&rank_type=Heat", nil)
	router.ServeHTTP(w, req)
}

func TestGeneralRouter(t *testing.T){
	cache.OfflineCacheClear()
	router:=RouterSet()
	w := httptest.NewRecorder()

	claims := cookie.MyClaims{
		"Admin",
		1,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 3600),
			Issuer:    "Gif-Dio",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(utils.COOKIE_SALT))
	
	req, _ := http.NewRequest(http.MethodPost, "/verify", nil)
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:tokenString,
	})
	router.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodPost, "/remove", nil)
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:tokenString,
	})
	router.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodGet, "/recommend", nil)
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:tokenString,
	})
	router.ServeHTTP(w, req)
	req, _ = http.NewRequest(http.MethodGet, "/favor", nil)
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:tokenString,
	})
	router.ServeHTTP(w, req)

}
