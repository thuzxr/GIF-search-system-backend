package register

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"backend/database"

	// "github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
)

func Register(c *gin.Context, db *sql.DB) string {
	username := c.DefaultQuery("user", "")
	password := c.DefaultQuery("password", "")
	// veri_input := c.DefaultQuery("vericode", "")
	// if !captcha.VerifyString(captchaId, veri_input) {
	// return "验证码错误"
	// }
	return database.InsertUser(username, password, "", db)
}
