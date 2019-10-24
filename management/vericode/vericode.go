package vericode

import (
	"bytes"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func get_vericode(c *gin.Context) string {
	length := captcha.DefaultLen
	captchaId := captcha.NewLen(length)
	w := c.Writer
	r := c.Request
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "10")

	var content bytes.Buffer
	w.Header().Set("Content-Type", "image/png")
	captcha.WriteImage(&content, captchaId, captcha.StdWidth, captcha.StdHeight)
	http.ServeContent(w, r, captchaId, time.Time{}, bytes.NewReader(content.Bytes()))
	return captchaId
}
