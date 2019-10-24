package recommend

import (
	"backend/utils"
)

func Recommend(gif utils.Gifs, gifs []utils.Gifs) []utils.Gifs {
	var recommend []utils.Gifs
	for _, idx := range gif.Recommend {
		recommend = append(recommend, gifs[idx])
	}
	return recommend
}
