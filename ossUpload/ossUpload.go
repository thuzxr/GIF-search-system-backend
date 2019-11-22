package ossUpload

import (
	"backend/utils"
	"fmt"
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//上传对应的Gif文件。一次性执行即可。localpath为gif存储文件夹相对路径
func OssUpload(gif utils.Gifs, localpath string) {
	log.SetPrefix("Error:")

	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "LTAI4FduW6Yf6AZY8ysPGmB9", "2eayaXUYwzCzK8HuOv8yrqRvtmsxd9")
	if err != nil {
		log.Panicln(err)
	}

	bucket, err := client.Bucket("gif-dio")
	if err != nil {
		log.Panicln(err)
	}

	err = bucket.PutObjectFromFile(gif.Name+".gif", localpath+"\\"+gif.Name+".gif")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(gif.Title + "Uploaded")
}

func OssMove(name string){
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "LTAI4FduW6Yf6AZY8ysPGmB9", "2eayaXUYwzCzK8HuOv8yrqRvtmsxd9")
	if err != nil {
		log.Panicln(err)
	}

	bucket, err := client.Bucket("gif-pre")
	if err != nil {
		log.Panicln(err)
	}

	bucket.CopyObjectTo("gif-dio", name, name)
	if err != nil {
		fmt.Println("CopyObjectTo Error:", err)
		log.Panicln(err)
	}

	// err = bucket.DeleteObject(name)
    // if err != nil {
    //     fmt.Println("Error:", err)
    //     log.Panicln(err)
    // }

}