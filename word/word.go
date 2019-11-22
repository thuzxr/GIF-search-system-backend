package word

import(
	"fmt"
	"github.com/go-ego/gse"
	"os"
	"io/ioutil"
	"encoding/json"
	"backend/utils"
	"strings"
	"sort"
	"time"
	"math"
)

func VecParse() map[string][][]float32{
	f, _:=ioutil.ReadFile("emb.json")
	m:=make(map[string][][]float32)
	json.Unmarshal(f,&m)
	return m
}

func DataCheck() bool{
	_,err:=os.Stat("emb_short.json")
	if os.IsNotExist(err){
		return false
	}
	_,err=os.Stat("embVectors.json")
	if os.IsNotExist(err){
		return false
	}
	return true
}

func FastVecParse() map[string][][]uint8{
	f, _:=ioutil.ReadFile("emb_short.json")
	m:=make(map[string][][]uint8)
	json.Unmarshal(f,&m)
	return m
}

func HammingCode(vec []uint8) []uint64{
	res:=[]uint64{0,0,0,0,0}
	for i:=0;i<300;i++{
		if(vec[i]>utils.HAMMING_EDGE){
		res[i/60]=res[i/60]|(1<<(i%60))
		}
	}
	return res
}

func HammingJudge(vec_1 []uint64, vec_2 []uint64, HAM_EDGE uint64) bool{
	// var EDGE uint64
	// EDGE=200
	// 50: 17/311 60: 25/311 100: 215/311
	var cnt uint64
	cnt=0
	for i:=0;i<5;i++{
		res:=vec_1[i]^vec_2[i]
		for j:=0;j<60;j++{
			cnt+=1-(res&1);
			res=res>>1;
		}
	}
	return cnt>HAM_EDGE
}

func HammingScreen(hamVec0 []uint64, hamVec [][]uint64, names []string, HAM_EDGE uint64) []string{
	res:=make([]string,0)
	for i:=range(hamVec){
		if(i>0){
			if strings.Compare(names[i-1], names[i])!=0 {
			continue
			}
		}
		if HammingJudge(hamVec0, hamVec[i], HAM_EDGE){
			res=append(res, names[i])
		}
	}
	return res
}

func WordToVecInit() map[string][]uint8{
	fmt.Println("WordToVec Initing")
	f, _:=ioutil.ReadFile("embVectors.json")
	m:=make(map[string][]uint8)
	json.Unmarshal(f,&m)
	fmt.Println("WordToVec Inited")
	return m
}

func WordToVec(keyword string, seg gse.Segmenter, m map[string][]uint8) [][]uint8{
	words:=gse.ToSlice(seg.Segment([]byte(keyword)), false)
	fmt.Println(words)
	// fmt.Println("vec_keyword:")
	res:=make([][]uint8,0)
	for i:=range(words){
		r0, b:=m[words[i]]
		if(b){
			res=append(res,r0)
		}
	}
	// fmt.Println(res)
	if(len(res)==0){
		res=append(res, make([]uint8, 300))
	}
	return res
}

func RankSearchInit() ([]string, map[string][][]uint8, [][]uint64){
	m:=FastVecParse()
	re_idx:=make([]string,0)
	// vecs:=make([][]uint8,0)
	vec_h:=make([][]uint64,0)
	for k, v:=range(m){
		for i:=range(v){
			re_idx=append(re_idx,k)
			vec_h=append(vec_h, HammingCode(v[i]))
			// vecs=append(vecs,v[i])
		}
	}
	fmt.Println("vec_HammingCode initiated")
	return re_idx, m, vec_h
}

func cosine(vec_1 []uint8, vec_2 []uint8) uint64{
	var ans uint64
	ans=0;
	var mod uint64
	for i:=0;i<300;i++{
		ans+=uint64(vec_1[i]+vec_2[i])
		mod+=uint64(vec_2[i]*vec_2[i])
	}
	if(mod<10){
		return 0;
	}else{		
		ans=ans/uint64(math.Sqrt(float64(mod)))
	}
	return ans
}

func Simple_Sim(vec_1 [][]uint8, vec_2 [][]uint8) uint64{
	var t,t0 uint64
	t=0
	for i:=range(vec_1){
		for j:=range(vec_2){
			t0=cosine(vec_2[j],vec_1[i])
			if(t0>t){
				t=t0;
			}
		}
	}
	return t
}

type sortRank struct{
	name string
	rank uint64
}

type sortRanks []sortRank

func (a sortRanks) Len() int {
	return len(a)
}
func (a sortRanks) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a sortRanks) Less(i, j int) bool { // 从大到小排序
	return a[j].rank < a[i].rank
}

func RankSearch(keyword string, word2vec map[string][]uint8, gif2vec map[string][][]uint8, 
	vec_h[][]uint64, re_idx []string, seg gse.Segmenter, HAM_EDGE uint64) []string{
	fmt.Println("searching")
	time_start:=time.Now()
	vec_keyword:=WordToVec(keyword, seg, word2vec)
	vec_0:=HammingCode(vec_keyword[0])
	pre_res:=HammingScreen(vec_0, vec_h, re_idx, HAM_EDGE)
	ranks:=make([]sortRank, len(pre_res))
	for i:=range(pre_res){
		ranks[i].name=pre_res[i]
		ranks[i].rank=Simple_Sim(gif2vec[pre_res[i]], vec_keyword)
	}
	sort.Sort(sortRanks(ranks))
	res2:=make([]string, len(pre_res))
	for i:=range(pre_res){
		res2[i]=ranks[i].name
	}
	time_2:=time.Since(time_start)
	fmt.Println(time_2)
	return res2
}

func Name_reIdx(gifs []utils.Gifs) map[string]*utils.Gifs{
	m:=make(map[string]*utils.Gifs)
	for i:=range(gifs){
		m[gifs[i].Name]=&gifs[i]
	}
	return m
}

func GifToVec(gif utils.Gifs, seg gse.Segmenter,m map[string][]uint8) ([][]uint8, [][]uint64, []string){
	re_idx:=make([]string, 0)
	veci:=WordToVec(gif.Keyword, seg, m);
	vec_h:=make([][]uint64, 0)
	for i:=range(veci){
		vec_h=append(vec_h, HammingCode(veci[i]))
		re_idx=append(re_idx, gif.Name)
	}
	return veci, vec_h, re_idx
}