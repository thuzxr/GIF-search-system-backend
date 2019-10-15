package main

import (
	db "backend/database"
	// "backend/utils"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// func SearchDemo(searchKey string, gifs []utils.Gifs) string {
// 	for i := range gifs {
// 		for j := range gifs[i].Keyword {
// 			if strings.Compare(gifs[i].Keyword[j], searchKey) == 0 {
// 				return ossUpload.OssSignLink(gifs[i], 3600)
// 			}
// 		}
// 	}
// 	return "Null"
// }

func RouterSet(DB *sql.DB) *gin.Engine {
	r := gin.Default()
	// gif := utils.JsonParse(".")
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
		c.JSON(200, gin.H{
			"message": "hello world! --sent by GO",
		})
	})
	r.GET("/search", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition")
		// searchKey := c.Query("key")
		// res := SearchDemo(searchKey, gif)
		keywords := c.DefaultQuery("key", "UNK")
		match := db.Query(DB, keywords) //search(keywords,gifs,DB)
		if len(match) == 0 {
			c.JSON(200, gin.H{
				"status": "failed",
			})
		} else {
			c.JSON(200, gin.H{
				"status": "succeed",
				"result": match,
			})
		}

	})
	// r.Run(":8000")
	return r
}

func main() {
	// gifs := utils.JsonParse(".")
	DB := db.Connect_db()
	// db.CreateTable(DB)
	// db.DB_init(gifs, DB)

	r := RouterSet(DB)
	r.Run(":80")

	// gifs:=ocr.JsonParse(".")
	// var gif []ocr.Gifs
	// fmt.Println(gif[0])
	// fmt.Println(SearchDemo("吐出来",gif))
}
