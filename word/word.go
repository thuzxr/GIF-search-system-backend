package word

import(
	"fmt"
	pinyin "github.com/mozillazg/go-pinyin"
)

func ConvertToPinyin(input string){
	args:=pinyin.NewArgs()
	fmt.Println(pinyin.Pinyin(input,args))
}