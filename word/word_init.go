package word

import(
	"backend/utils"
	"encoding/json"
	"os"
)

func EmbZip(){
	m:=VecParse()
	m0:=make(map[string][][]uint8)
	var intLis []uint8
	var intLis0 [][]uint8
	for key, value :=range(m){
		intLis0=make([][]uint8,0)
		for i:=0;i<len(value);i++{
			intLis=make([]uint8,300)
			for j:=0;j<300;j++{
				intLis[j]=uint8(utils.HAMMING_EDGE+utils.HAMMING_DIV*value[i][j])
			}
			intLis0=append(intLis0,intLis)
		}
		m0[key]=intLis0
	}
	// fmt.Println(m0)
	b, _:=json.Marshal(m0)
	w1, _ := os.OpenFile("emb_short.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	_, _ = w1.Write(b)
	_ = w1.Close()
}