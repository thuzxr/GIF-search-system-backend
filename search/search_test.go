package search

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"backend/utils"
)

func TestSearch(t *testing.T) {
	os.Chdir("..")
	names, titles, keywords := FastIndexParse()
	gifs:=make([]utils.Gifs, 0)
	for i:=0;i<len(names);i++{
		gifs=append(gifs, utils.Gifs{
			Name:names[i],
			Title:titles[i],
			Keyword:keywords[i],
		})
	}
	fmt.Println(len(names), len(titles), len(keywords))
	keyword := "开心"
	match := SimpleSearch(keyword, gifs)
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
