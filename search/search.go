package search

import (
	"backend/utils"
	"strings"
)

//简单的离线搜索算法
func SimpleSearch(keyword string, gifs []utils.Gifs) []utils.Gifs {
	var ans []utils.Gifs
	for i := 0; i < len(gifs); i++ {
		if strings.Contains(gifs[i].Keyword, keyword) || strings.Contains(gifs[i].Title, keyword) {
			ans = append(ans, gifs[i])
		}
	}
	return ans
}
