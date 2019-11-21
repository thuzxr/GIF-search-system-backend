package login

import (
	"backend/database"
	"database/sql"
)

func Login(user, password string, DB *sql.DB) int {
	return database.QueryUser(user, password, DB)
}
