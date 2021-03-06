package database

import (
	"database/sql"
	"fmt"
	"time"

	"backend/utils"

	_ "github.com/go-sql-driver/mysql"

	goini "github.com/clod-moon/goconf"
)

func ConnectDB(path string) *sql.DB {
	conf := goini.InitConfig(path)
	serverAddr := conf.GetValue("database", "server")
	fmt.Println(serverAddr)
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", utils.USERNAME, utils.PASSWORD, utils.NETWORK, serverAddr, utils.PORT, utils.DATABASE)
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

func ErrProc(err error) bool {
	if err != nil {
		fmt.Println("create table failed:", err)
		return false
	}
	return true
}

func Init(DB *sql.DB) {

	sql := `CREATE TABLE IF NOT EXISTS USER_MANAGE(
			USER  		VARCHAR(64) 	NOT NULL 	PRIMARY KEY ,
			PASSWORD	TEXT			NOT NULL,
			TYPE		INTEGER			DEFAULT 1
			);`
	_, err := DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table USER_MANAGE succeed")

	sql = `CREATE TABLE IF NOT EXISTS PROFILE(
		USER  		VARCHAR(64) 	NOT NULL 	PRIMARY KEY ,
		Email		TEXT ,
		FirstName	TEXT ,
		LastName	TEXT ,
		Addr		TEXT ,
		ZipCode		TEXT ,
		City 		TEXT ,
		Country		TEXT ,
		Birthday 	TEXT ,
		Height		TEXT ,
		About		TEXT ,
		FOREIGN KEY(USER) REFERENCES USER_MANAGE(USER) ON DELETE CASCADE
		);`
	_, err = DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table PROFILE succeed")

	//Follow table
	sql = `CREATE TABLE IF NOT EXISTS FOLLOW(
		USER  		VARCHAR(64) 	NOT NULL,
		Follows 	VARCHAR(64) 	NOT NULL,
		FOREIGN KEY(USER) REFERENCES PROFILE(USER) ON DELETE CASCADE,
		PRIMARY KEY(USER, Follows)
		);`
	_, err = DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table FOLLOW succeed")

	sql = `CREATE TABLE IF NOT EXISTS FAVOR(
		USER	VARCHAR(64)			NOT NULL,
		GifId	VARCHAR(64)			NOT NULL,
		PRIMARY KEY(USER, GifId),
		FOREIGN KEY(USER) REFERENCES PROFILE(USER) ON DELETE CASCADE
		);`
	_, err = DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table FAVOR succeed")

	sql = `CREATE TABLE IF NOT EXISTS GIF_INFO(
		USER	VARCHAR(64)		NOT NULL,
		GifId	VARCHAR(64)		NOT NULL	PRIMARY KEY,
		TAG		TEXT,
		INFO	TEXT,
		TITLE	TEXT,
		FOREIGN KEY(USER) REFERENCES PROFILE(USER) ON DELETE CASCADE
	);`
	_, err = DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table GIF_INFO succeed")

	sql = `CREATE TABLE IF NOT EXISTS GIF_TOVERIFY(
		USER	VARCHAR(64)		NOT NULL,
		GifId	VARCHAR(64)		NOT NULL	PRIMARY KEY,
		TAG		TEXT,
		INFO	TEXT,
		TITLE	TEXT,
		FOREIGN KEY(USER) REFERENCES PROFILE(USER) ON DELETE CASCADE
	);`
	_, err = DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table GIF_TOVERIFY succeed")

	//Comments table
	sql = `CREATE TABLE IF NOT EXISTS COMMENTS(
		ComId 	INTEGER 	NOT NULL,
		GifId 	VARCHAR(64) NOT NULL,
		Comment TEXT,
		User 	VARCHAR(64),
		PRIMARY KEY (ComId, GifId),
		FOREIGN KEY(GifId) REFERENCES GIF_INFO(GifId) ON DELETE CASCADE
		); `
	_, err = DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table COMMENTS succeed")

	//like table
	sql = `CREATE TABLE IF NOT EXISTS LIKES(
		GifId 	   VARCHAR(64)  NOT NULL,
		USER       VARCHAR(64)  NOT NULL,
		PRIMARY KEY(GifId,USER),
		FOREIGN KEY(GifId) REFERENCES GIF_INFO(GifId) ON DELETE CASCADE
		); `
	_, err = DB.Exec(sql)
	if !ErrProc(err) {
		return
	}
	fmt.Println("create table LIKES succeed")
}

func InsertUser(user, password string, DB *sql.DB) string {
	sql := `select USER from USER_MANAGE where USER='` + user + `'`
	rows, _ := DB.Query(sql)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if rows.Next() {
		return "用户名已存在"
	}

	sql = `insert INTO USER_MANAGE(USER,PASSWORD) values('` + user + `','` + password + `')`
	_, err := DB.Exec(sql)
	if err != nil {
		print(err)
	}

	_, err = DB.Exec(`insert INTO PROFILE(USER) values(?)`, user)
	if err != nil {
		print(err)
	}
	return "注册成功！"
}

func InsertGIF(DB *sql.DB, user, GifId, TAG, INFO, TITLE string) {
	_, err := DB.Exec("insert INTO GIF_INFO(USER,GifId,TAG,INFO,TITLE) values(?,?,?,?,?)", user, GifId, TAG, INFO, TITLE)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		// fmt.Print(gif)
		return
	}
}

func InsertUnderVerifyGIF(DB *sql.DB, user, GifId, TAG, INFO, TITLE string) {
	_, err := DB.Exec("insert INTO GIF_TOVERIFY(USER,GifId,TAG,INFO,TITLE) values(?,?,?,?,?)", user, GifId, TAG, INFO, TITLE)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}
}

func GetToVerifyGIF(DB *sql.DB) []QueryGif {
	rows, qerr := DB.Query("select USER,GifId,TAG,INFO,TITLE from GIF_TOVERIFY")

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if qerr != nil {
		fmt.Printf("query failed, err:%v\n", qerr)
	}
	res := make([]QueryGif, 0)
	var user, gifId, tag, info, title string
	for rows.Next() {
		if serr := rows.Scan(&user, &gifId, &tag, &info, &title); serr != nil {
			fmt.Printf(utils.ScanFailed, serr)
		}
		res = append(res, QueryGif{
			GifId: gifId,
			TAG:   tag,
			INFO:  info,
			TITLE: title,
		})
	}

	return res
}

func VerifyGIF(DB *sql.DB, GifId string) {
	rows, qerr := DB.Query("select USER,GifId,TAG,INFO,TITLE from GIF_TOVERIFY WHERE GifId like '%" + GifId + "%'")

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if qerr != nil {
		fmt.Printf("query failed, err:%v\n", qerr)
	}
	var user, gifId, tag, info, title string
	for rows.Next() {
		if serr := rows.Scan(&user, &gifId, &tag, &info, &title); serr != nil {
			fmt.Printf(utils.ScanFailed, serr)
		}
		InsertGIF(DB, user, gifId, tag, info, title)
		_, err := DB.Exec(`DELETE FROM GIF_TOVERIFY WHERE GifId='` + GifId + `'`)
		if err != nil {
			fmt.Println("error in delete gif from toverify %v", err)
		}
	}
}

func RemoveVerify(DB *sql.DB, GifId string) {
	_, err := DB.Exec(`DELETE FROM GIF_TOVERIFY WHERE GifId='` + GifId + `'`)
	if err != nil {
		fmt.Println("error in delete gif from toverify %v", err)
	}
}

// func InsertFavor(user, GifId string, DB *sql.DB) string {
// 	_, err := DB.Exec("INSERT INTO FAVOR(USER,GifId) values(?,?)", user, GifId)
// 	if err != nil {
// 		return "收藏失败"
// 	} else {
// 		return "收藏成功"
// 	}
// }

func UpdateLikes(likes map[string][]string, DB *sql.DB) {
	DB.Exec("DELETE * from LIKES")
	if len(likes) != 0 {
		sql := "INSERT INTO LIKES(GifId,USER) VALUES"
		all_str := ""
		for gifid := range likes {
			for i := range likes[gifid] {
				all_str += ",('" + gifid + "','" + likes[gifid][i] + "')"
			}
		}
		sql += all_str[1:]

		_, err := DB.Exec(sql)
		if err != nil {
			fmt.Println("失败")
		} else {
			fmt.Println("成功")
		}
	}
	//sql := "INSERT INTO LIKES(GifId,USER) VALUES"
	//all_str := ""
	//for gifid := range likes {
	//	for i := range likes[gifid] {
	//		all_str += ",('" + gifid + "','" + likes[gifid][i] + "')"
	//	}
	//}
	//sql += all_str[1:]

	//_, err := DB.Exec(sql)
	//if err != nil {
	//	fmt.Println("失败")
	//} else {
	//	fmt.Println("成功")
	//}
}

// func InsertFollow(user, follow string, DB *sql.DB) string {
// 	_, err := DB.Exec("INSERT INTO FOLLOW(USER,Follows) values(?,?)", user, follow)
// 	if err != nil {
// 		return "关注失败"
// 	} else {
// 		return "关注成功"
// 	}
// }

func ChangeProfile(user string, Profile utils.Profile, DB *sql.DB) string {
	_, err := DB.Exec(`UPDATE PROFILE SET Email='` + Profile.Email + `', FirstName='` + Profile.FirstName + `', LastName='` + Profile.LastName + `', Addr='` + Profile.Addr + `',ZipCode='` + Profile.ZipCode + `', City='` + Profile.City + `', Country='` + Profile.Country + `', About='` + Profile.About + `', Height='` + Profile.Height + `', Birthday='` + Profile.Birthday + `' WHERE USER='` + user + `'`)
	if err != nil {
		return "更新失败"
	} else {
		return "更新成功"
	}
}

// func DeleteFavor(user string, GifIds []string, DB *sql.DB) string {
// 	var GifId string
// 	for gifid := range GifIds {
// 		gifname := "'" + GifIds[gifid] + "'"
// 		GifId = GifId + "," + gifname
// 	}
// 	GifId = GifId[1:]
// 	_, err := DB.Exec(`DELETE FROM FAVOR WHERE USER='` + user + `' AND GifId IN(` + GifId + `)`)
// 	if err != nil {
// 		return "删除错误"
// 	} else {
// 		return "删除成功"
// 	}
// }

// func DeleteFollow(user, follow string, DB *sql.DB) string {
// 	_, err := DB.Exec(`DELETE FROM FOLLOW WHERE USER='` + user + `' AND Follows='` + follow + `'`)
// 	if err != nil {
// 		return "删除关注失败"
// 	} else {
// 		return "删除关注成功"
// 	}
// }

// func DeleteComment(commentId, GifId string, DB *sql.DB) string {
// 	_, err := DB.Exec(`DELETE FROM COMMENTS WHERE commentId=` + commentId + ` AND GifId='` + GifId + `'`)
// 	if err != nil {
// 		return "删除评论失败"
// 	} else {
// 		return "删除评论成功"
// 	}
// }

func DeleteGif(GifId string, DB *sql.DB) string {
	_, err := DB.Exec(`DELETE FROM GIF_INFO WHERE GifId='` + GifId + `'`)
	if err != nil {
		return "删除图片失败"
	} else {
		return "删除图片成功"
	}
}

func DeleteAccount(user string, DB *sql.DB) string {
	_, err := DB.Exec(`DELETE FROM USER_MANAGE WHERE USER='` + user + `'`)
	if err != nil {
		return "注销失败"
	} else {
		return "注销成功"
	}
}

func QueryProfile(user string, DB *sql.DB) []string {
	var Email, FirstName, LastName, Addr, ZipCode, City, Country, About, Height, Birthday string
	rows, _ := DB.Query("select Email, FirstName, LastName, Addr, ZipCode, City, Country, About,Height,Birthday from PROFILE WHERE USER='" + user + "'")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	rows.Next()
	err := rows.Scan(&Email, &FirstName, &LastName, &Addr, &ZipCode, &City, &Country, &About, &Height, &Birthday)
	if err != nil {
		fmt.Println("查询失败！")
	}
	var returns []string
	returns = append(returns, Email, FirstName, LastName, Addr, ZipCode, City, Country, About, Height, Birthday)
	return returns
}

// func QueryFavor(user string, DB *sql.DB) []string {
// 	var favors []string
// 	var favor string
// 	rows, _ := DB.Query("select GifId from FAVOR WHERE USER='" + user + "'")
// 	defer func() {
// 		if rows != nil {
// 			rows.Close()
// 		}
// 	}()
// 	for rows.Next() {
// 		err := rows.Scan(&favor)
// 		if err != nil {
// 			fmt.Printf("scan failed, err:%v\n", err)
// 		}
// 		favors = append(favors, favor)
// 	}
// 	return favors
// }

// func QueryFollow(user string, DB *sql.DB) []string {
// 	var follows []string
// 	var follow string
// 	rows, _ := DB.Query("select Follows from FOLLOW WHERE USER='" + user + "'")
// 	defer func() {
// 		if rows != nil {
// 			rows.Close()
// 		}
// 	}()
// 	for rows.Next() {
// 		err := rows.Scan(&follow)
// 		if err != nil {
// 			fmt.Printf("scan failed, err:%v\n", err)
// 		}
// 		follows = append(follows, follow)
// 	}
// 	return follows
// }

// func QueryFollower(user string, DB *sql.DB) []string {
// 	var followers []string
// 	var follower string
// 	rows, _ := DB.Query("select USER from FOLLOW WHERE Follows='" + user + "'")
// 	defer func() {
// 		if rows != nil {
// 			rows.Close()
// 		}
// 	}()
// 	for rows.Next() {
// 		err := rows.Scan(&follower)
// 		if err != nil {
// 			fmt.Printf("scan failed, err:%v\n", err)
// 		}
// 		followers = append(followers, follower)
// 	}
// 	return followers
// }

// type Comment struct {
// 	ComId   int
// 	Comment string
// }

// func QueryComment(GifId string, DB *sql.DB) []Comment {
// 	var comments []Comment
// 	comment := new(Comment)
// 	rows, _ := DB.Query("select ComId, Comment from COMMENTS WHERE GifId='" + GifId + "'")
// 	defer func() {
// 		if rows != nil {
// 			rows.Close()
// 		}
// 	}()
// 	for rows.Next() {
// 		err := rows.Scan(&comment.ComId, &comment.Comment)
// 		if err != nil {
// 			fmt.Printf("scan failed, err:%v\n", err)
// 		}
// 		comments = append(comments, *comment)
// 	}
// 	return comments
// }

type QueryGif struct {
	GifId  string
	TAG    string
	INFO   string
	TITLE  string
	OSSURL string
}

func QueryGifs(user string, DB *sql.DB) []QueryGif {
	var QGifs []QueryGif
	querygif := new(QueryGif)
	rows, _ := DB.Query("select GifId,TAG,INFO,TITLE from GIF_INFO WHERE USER='" + user + "'")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	for rows.Next() {
		err := rows.Scan(&querygif.GifId, &querygif.TAG, &querygif.INFO, &querygif.TITLE)
		if err != nil {
			fmt.Printf(utils.ScanFailed, err)
		}
		QGifs = append(QGifs, *querygif)
	}
	return QGifs
}

func QueryUser(user, password string, DB *sql.DB) int {
	var user_type int
	rows, _ := DB.Query("select TYPE from USER_MANAGE WHERE USER='" + user + "' AND PASSWORD='" + password + "'")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if rows.Next() {
		rows.Scan(&user_type)
	} else {
		user_type = -1
	}
	return user_type
}

func LoadAll(DB *sql.DB) ([]string, []string, []utils.Gifs, map[string][]string, map[string][]string) {
	var users []string
	// var names []string    //id
	// var titles []string   //title
	var infos []string //info
	// var keywords []string //tags
	gifs := make([]utils.Gifs, 0)

	var user string
	var name string
	var title string
	var info string
	var keyword string

	rows, _ := DB.Query("Select USER,GifId,TAG,INFO,TITLE FROM GIF_INFO")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	for rows.Next() {
		rows.Scan(&user, &name, &keyword, &info, &title)
		users = append(users, user)
		infos = append(infos, info)
		// names = append(names, name)
		// titles = append(titles, title)
		// keywords = append(keywords, keyword)
		gifs = append(gifs, utils.Gifs{
			Name:    name,
			Title:   title,
			Keyword: keyword,
		})
	}
	if rows != nil {
		rows.Close()
	}

	var gifid string

	var likes map[string][]string
	likes = make(map[string][]string)

	var likes_u2g map[string][]string
	likes_u2g = make(map[string][]string)

	rows, _ = DB.Query("Select GifId,USER FROM LIKES")
	for rows.Next() {
		rows.Scan(&gifid, &user)
		_, ok := likes[gifid]
		if !ok {
			likes[gifid] = make([]string, 0)
		}
		likes[gifid] = append(likes[gifid], user)

		_, ok = likes_u2g[user]
		if !ok {
			likes_u2g[user] = make([]string, 0)
		}
		likes_u2g[user] = append(likes_u2g[gifid], user)
	}

	return users, infos, gifs, likes, likes_u2g
}
