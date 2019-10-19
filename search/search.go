package search

import (
	"backend/utils"
	"strings"
)

//简单的离线搜索算法
func SimpleSearch(keyword string, names []string, titles []string, keywords []string) []utils.Gifs {
	var ans []utils.Gifs
	gif := new(utils.Gifs)
	for i := 0; i < len(names); i++ {
		if strings.Contains(keywords[i], keyword) {
			gif.Name = names[i]
			gif.Keyword = keywords[i]
			gif.Title = titles[i]
			ans = append(ans, *gif)
		}
	}
	return ans
}
