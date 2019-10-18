package main

import (
	"fmt"

	"math/rand"
	"time"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	model, err := tf.LoadSavedModel("/Users/saberrrrrrrr/go/src/tf_test/python_models/models/CBOW", []string{"var"}, nil)

	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return
	}

	defer model.Session.Close()

	var array1 [1][]int32
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		array1[0] = append(array1[0], seed.Int31n(2000)) //生成随机数
	}

	var array2 [1][]int32
	seed2 := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		array2[0] = append(array2[0], seed2.Int31n(2000)) //生成随机数
	}

	tensor1, terr := tf.NewTensor(array1)
	if terr != nil {
		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
		return
	}

	tensor2, terr := tf.NewTensor(array2)
	if terr != nil {
		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
		return
	}

	result, runErr := model.Session.Run(
		map[tf.Output]*tf.Tensor{
			model.Graph.Operation("CBOW/placeholder/input1").Output(0): tensor1,
			model.Graph.Operation("CBOW/placeholder/input2").Output(0): tensor2,
		},
		[]tf.Output{
			model.Graph.Operation("CBOW/Sim/cossim").Output(0),
		},
		nil,
	)

	if runErr != nil {
		fmt.Printf("Error running the session with input, err: %s\n", runErr.Error())
		return
	}

	fmt.Printf("%v", result[0].Value())
}
