package cbow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"backend/utils"

	"github.com/go-ego/gse"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func graph_load(path string) *tf.SavedModel {
	model, err := tf.LoadSavedModel(path, []string{"var"}, nil)
	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return model
	}
	return model
}

func seg_init() gse.Segmenter {
	var seg gse.Segmenter
	seg.LoadDict()
	return seg
}

func load_word2idx(path string) map[string]int32 {
	word2idx := make(map[string]int32)
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error loading file: %s\n", err.Error())
		return word2idx
	}
	text, err := ioutil.ReadAll(f)
	json.Unmarshal(text, &word2idx)
	return word2idx
}

func get_word(tag []string, word2idx map[string]int32, pad_length int) []int32 {
	var idx []int32
	for _, word := range tag {
		id, exist := word2idx[word]
		if exist {
			idx = append(idx, id)
		} else {
			idx = append(idx, 1)
		}
	}
	for i := len(idx); i < pad_length; i++ {
		idx = append(idx, 0)
	}
	return idx
}

func preprocess(tags []string, seg gse.Segmenter, word2idx map[string]int32) [][]int32 {
	var res [][]int32
	for _, tag := range tags {
		tag_byte := []byte(tag)
		splits := seg.Segment(tag_byte)
		word_idx := get_word(gse.ToSlice(splits), word2idx, 20)
		res = append(res, word_idx)

	}
	return res
}

func vector_map(tags [][]int32, model *tf.SavedModel) *tf.Tensor {
	tensor, terr := tf.NewTensor(tags)
	if terr != nil {
		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
		return nil
	}

	result, runErr := model.Session.Run(
		map[tf.Output]*tf.Tensor{
			model.Graph.Operation("CBOW/placeholder/input1").Output(0): tensor,
		},
		[]tf.Output{
			model.Graph.Operation("CBOW/Emb/emb1").Output(0),
		},
		nil,
	)

	if runErr != nil {
		fmt.Printf("Error running the session with input, err: %s\n", runErr.Error())
		return nil
	}
	return result[0]
}

func compute_sim(tags1 [][]int32, tags2 [][]int32, model *tf.SavedModel) *tf.Tensor {
	tensor1, terr := tf.NewTensor(tags1)
	if terr != nil {
		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
		return nil
	}

	tensor2, terr := tf.NewTensor(tags2)
	if terr != nil {
		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
		return nil
	}
	result, runErr := model.Session.Run(
		map[tf.Output]*tf.Tensor{
			model.Graph.Operation("CBOW/placeholder/input1").Output(0): tensor1,
			model.Graph.Operation("CBOW/placeholder/input2").Output(0): tensor2,
		},
		[]tf.Output{
			model.Graph.Operation("CBOW/Sim/mean_sim").Output(0),
		},
		nil,
	)

	if runErr != nil {
		fmt.Printf("Error running the session with input, err: %s\n", runErr.Error())
		return nil
	}
	return result[0]
}

func Init(path_model string, path_word2idx string, gifs []utils.Gifs) *tf.SavedModel {
	model := graph_load(path_model)
	seg := seg_init()
	word2idx := load_word2idx(path_word2idx)
	for i := 0; i < len(gifs); i++ {
		gifs[i].Word_idx = preprocess(strings.Split(gifs[i].Keyword, " "), seg, word2idx)
	}
	return model
}

type sortgif struct {
	index int
	score float32
}

type sortgifslice []sortgif

func (a sortgifslice) Len() int {
	return len(a)
}
func (a sortgifslice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a sortgifslice) Less(i, j int) bool { // 从大到小排序
	return a[j].score < a[i].score
}

func Recommend(gif utils.Gifs, gifs []utils.Gifs, model *tf.SavedModel) []utils.Gifs {
	var scores []sortgif
	for i, gif2 := range gifs {
		score := compute_sim(gif.Word_idx, gif2.Word_idx, model).Value().(float32)
		var append_gif sortgif
		append_gif.score = score
		append_gif.index = i
		scores = append(scores, append_gif)
	}
	sort.Sort(sortgifslice(scores))
	var commend_gifs []utils.Gifs
	for i := 0; i < 10; i++ {
		commend_gifs = append(commend_gifs, gifs[scores[i].index])
	}
	return commend_gifs
}

// func Test(path_model string, path_word2idx string, gif1 utils.Gifs, gif2 utils.Gifs) interface{} {
// 	model := graph_load(path_model)
// 	seg := seg_init()
// 	word2idx := load_word2idx(path_word2idx)

// 	tags1 := strings.Split(gif1.Keyword, " ")
// 	tags2 := strings.Split(gif2.Keyword, " ")

// 	tags1_pre := preprocess(tags1, seg, word2idx)
// 	tags2_pre := preprocess(tags2, seg, word2idx)

// 	sim := compute_sim(tags1_pre, tags2_pre, model)
// 	ret := sim.Value()
// 	return ret
// }

// func upload() {
// 	model, err := tf.LoadSavedModel("/Users/saberrrrrrrr/go/src/tf_test/python_models/models/CBOW", []string{"var"}, nil)

// 	if err != nil {
// 		fmt.Printf("Error loading saved model: %s\n", err.Error())
// 		return
// 	}

// 	defer model.Session.Close()

// 	var array1 [1][]int32
// 	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	for i := 0; i < 100; i++ {
// 		array1[0] = append(array1[0], seed.Int31n(2000)) //生成随机数
// 	}

// 	var array2 [1][]int32
// 	seed2 := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	for i := 0; i < 100; i++ {
// 		array2[0] = append(array2[0], seed2.Int31n(2000)) //生成随机数
// 	}

// 	tensor1, terr := tf.NewTensor(array1)
// 	if terr != nil {
// 		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
// 		return
// 	}

// 	tensor2, terr := tf.NewTensor(array2)
// 	if terr != nil {
// 		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
// 		return
// 	}

// 	result, runErr := model.Session.Run(
// 		map[tf.Output]*tf.Tensor{
// 			model.Graph.Operation("CBOW/placeholder/input1").Output(0): tensor1,
// 			model.Graph.Operation("CBOW/placeholder/input2").Output(0): tensor2,
// 		},
// 		[]tf.Output{
// 			model.Graph.Operation("CBOW/Sim/cossim").Output(0),
// 		},
// 		nil,
// 	)

// 	if runErr != nil {
// 		fmt.Printf("Error running the session with input, err: %s\n", runErr.Error())
// 		return
// 	}

// 	fmt.Printf("%v", result[0].Value())
// }
