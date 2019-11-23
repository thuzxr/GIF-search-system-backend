package register

import(
	"backend/database"
	"testing"
	"github.com/gin-gonic/gin"

	"net/http"
	"net/http/httptest"
)

func TestRegister(t *testing.T){
	r:=gin.Default()
	DB:=database.ConnectDB()
	r.GET("/", func(c *gin.Context){
		c.String(200, Register(c, DB))
	})

	w:=httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
}