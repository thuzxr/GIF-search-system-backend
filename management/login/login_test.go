package login

import(
	"backend/database"
	"testing"
)

func TestLogin(t *testing.T){
	DB:=database.ConnectDB()
	Login("", "", DB)
}