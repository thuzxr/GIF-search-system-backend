package login

import (
	"backend/database"
	"database/sql"
)

func Login(user, password string, DB *sql.DB) string {
	return database.QueryUser(user, password, DB)
}
