package database

import (
	"database/sql"
	"fmt"
	// "time"

	// "backend/utils"

	_ "github.com/go-sql-driver/mysql"
)

func ErrProc(err error) bool{
	if err!=nil{
		fmt.Println("create table failed:" ,err)
		return false
	}
	return true
}

func GrandInit(DB *sql.DB) {
	//Profile table
	sql := `CREATE TABLE IF NOT EXISTS PROFILE(
		USER VARCHAR(64) PRIMARY KEY NOT NULL,
		Birthday TEXT
		); `
	_, err := DB.Exec(sql)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create table PROFILE succeed")

	//Follow table
	sql = `CREATE TABLE IF NOT EXISTS FOLLOW(
		USER VARCHAR(64) NOT NULL,
		Follows VARCHAR(64) NOT NULL,
		PRIMARY KEY (USER, Follows)
		); `
	_, err = DB.Exec(sql)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create table FOLLOW succeed")

	//Follower table
	sql = `CREATE TABLE IF NOT EXISTS FOLLOWER(
		USER VARCHAR(64) NOT NULL,
		Follower VARCHAR(64) NOT NULL,
		PRIMARY KEY (USER, Follower)
		); `
	_, err = DB.Exec(sql)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create table FOLLOWER succeed")

	//Comments table
	sql = `CREATE TABLE IF NOT EXISTS COMMENTS(
		Comment_id VARCHAR(64) NOT NULL,
		Gifs_id VARCHAR(64) NOT NULL,
		Comment TEXT,
		Time TEXT,
		User VARCHAR(64),
		PRIMARY KEY (Comment_id, Gifs_id)
		); `
	_, err = DB.Exec(sql)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create table COMMENTS succeed")
}

func ForeignInit(DB *sql.DB){
	_, err := DB.Exec(`alter table USER_INFO add constraint FK_user foreign key (USER) references PROFILE (USER)`)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create foreign from USER_INFO succeed")

	_, err = DB.Exec(`alter table PROFILE add constraint FK_user foreign key (USER) references FOLLOW (USER)`)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create foreign from PROFILE succeed")

	_, err = DB.Exec(`alter table PROFILE add constraint FK_user foreign key (USER) references FOLLOWER (USER)`)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create foreign from PROFILE succeed")
	
	_, err = DB.Exec(`alter table USER_INFO add constraint FK_user foreign key (USER) references PROFILE (USER)`)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create foreign from USER_INFO succeed")

	// _, err = DB.Exec(`alter table PROFILE add constraint FK_user foreign key (USER) references GIF_INTO (Username)`)
	// if ErrProc(err)==false{
	// 	return
	// }
	// fmt.Println("create foreign from PROFILE succeed")

	
	_, err = DB.Exec(`alter table GIF_INFO add constraint FK_user foreign key (Name) references COMMENTS (Gifs_id)`)
	if ErrProc(err)==false{
		return
	}
	fmt.Println("create foreign from USER_INFO succeed")
}