package login

import(
	"backend/database"
	"testing"
)

func TestLogin(t *testing.T){
	DB:=database.ConnectDB("../../settings.ini")
	Login("", "", DB)
}