package database

import (
	"database/sql"
	"fmt"
	"time"

	"backend/ossUpload"
	"backend/utils"

	_ "github.com/go-sql-driver/mysql"
)

func Connect_db() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", utils.USERNAME, utils.PASSWORD, utils.NETWORK, utils.SERVER, utils.PORT, utils.DATABASE)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Open mysql failed,err:%v\n", err)
		return DB
	}
	DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(100)                  //设置最大连接数
	DB.SetMaxIdleConns(16)                   //设置闲置连接数
	return DB
}

func DB_init(gifs []utils.Gifs, DB *sql.DB) {
	_, err := DB.Exec("DELETE FROM GIF_INFO")
	if err != nil {
		fmt.Println("delete data in table failed:", err)
		return
	}
	var idx int
	for idx = 0; idx < len(gifs); idx++ {
		InsertData(DB, gifs[idx])
	}
}

func CreateTable(DB *sql.DB) {
	sql := `CREATE TABLE IF NOT EXISTS GIF_INFO(
	Name VARCHAR(64) PRIMARY KEY NOT NULL,
	Keyword TEXT,
	Title TEXT,
	Gif_url TEXT,
	Cover_url TEXT,
	Oss_url TEXT
	); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}

	DB.Exec("alter table GIF_INFO convert to character set utf8mb4 collate utf8mb4_bin")

	fmt.Println("create table successd")
}

func InsertData(DB *sql.DB, gif utils.Gifs) {
	_, err := DB.Exec("insert INTO GIF_INFO(Name,Keyword,Title,Gif_url,Cover_url,Oss_url) values(?,?,?,?,?,?)", gif.Name, gif.Keyword, gif.Title, gif.Gif_url, gif.Cover_url, gif.Oss_url)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		// fmt.Print(gif)
		return
	}
	fmt.Println(gif.Title,"Inserted key:",gif.Keyword)
}

func Query(DB *sql.DB, q string) []utils.Gifs {
	gif := new(utils.Gifs)
	var ans []utils.Gifs
	fmt.Print(q)
	rows, qerr := DB.Query("select Name,Title,Keyword,Gif_url from GIF_INFO WHERE Keyword like '%" + q + "%'")

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if qerr != nil {
		fmt.Printf("query failed, err:%v\n", qerr)
		return ans
	}
	for rows.Next() {
		if serr := rows.Scan(&gif.Name,&gif.Title,&gif.Keyword,&gif.Gif_url); serr != nil {
			fmt.Printf("scan failed, err:%v\n", serr)
			return ans
		}
		gif.Oss_url = ossUpload.OssSignLink(*gif, 3600)
		ans = append(ans, *gif)
		fmt.Println(gif.Oss_url)
	}
	return ans
}
