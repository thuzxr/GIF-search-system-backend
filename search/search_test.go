package search

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	os.Chdir("..")
	names, titles, keywords := FastIndexParse()
	fmt.Println(len(names), len(titles), len(keywords))
	keyword := "开心"
	match := SimpleSearch(keyword, names, titles, keywords)
	assert.NotEqual(t, len(match), 0)
}
