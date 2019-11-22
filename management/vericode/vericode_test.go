package vericode

import(
	"net/http"
	"net/http/httptest"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVericode(t *testing.T){
	r:=gin.Default()
	r.GET("/", func(c *gin.Context){
		Get_vericode(c)
		Gen_vericode(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	captchaId := jsoniter.Get(w.Body.Bytes(), "captchaId").ToString()
	assert.NotEqual(t, captchaId, "")
}