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

func TestIndex(t *testing.T) {
	os.Chdir("..")
	gifs := IndexParse()
	assert.Equal(t, len(gifs), 0)
	names := NameIndex()
	assert.NotEqual(t, len(names), 0)
	titles := TitleIndex()
	assert.NotEqual(t, len(titles), 0)
	keyword := KeywordIndex()
	assert.NotEqual(t, len(keyword), 0)
}
