package main

import (
	"backend/cache"
	"net/http"
	"net/http/httptest"
	"testing"

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
}
