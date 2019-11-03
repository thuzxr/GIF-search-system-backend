package vericode

import (
	"bytes"
	"net/http"
	"path"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func Get_vericode(c *gin.Context) {
	length := captcha.DefaultLen
	captchaId := captcha.NewLen(length)
	c.JSON(http.StatusOK, gin.H{
		"captchaId": captchaId,
	})
}

func Gen_vericode(c *gin.Context) {
	w := c.Writer
	r := c.Request
	// captchaId := c.Param("captchaId")

	_, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	captchaId := file[:len(file)-len(ext)]
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "10")

	var content bytes.Buffer
	w.Header().Set("Content-Type", "image/png")
	captcha.WriteImage(&content, captchaId, captcha.StdWidth, captcha.StdHeight)
	http.ServeContent(w, r, captchaId+".png", time.Time{}, bytes.NewReader(content.Bytes()))
}
