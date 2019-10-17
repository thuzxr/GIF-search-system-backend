package word

import(
	"fmt"
	pinyin "github.com/mozillazg/go-pinyin"
)

func ConvertToPinyin(input string){
	args:=pinyin.NewArgs()
	lis0:=pinyin.Pinyin(input,args)
	ans:=make([]string,len(lis0))
	for i:=0;i<len(lis0);i++{
		ans[i]=lis0[i][0]
	}
	fmt.Println(ans)
}