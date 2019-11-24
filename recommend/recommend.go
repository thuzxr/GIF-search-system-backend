package recommend

import (
	"backend/utils"
	"math"
	"sort"
)

func Recommend(gif utils.Gifs, gifs []utils.Gifs) []utils.Gifs {
	var recommend []utils.Gifs
	for _, idx := range gif.Recommend {
		recommend = append(recommend, gifs[idx])
	}
	return recommend
}

// func

func intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v] = 1
	}

	for _, v := range slice2 {
		_, ok := m[v]
		if ok {
			nn = append(nn, v)
		}
	}
	return nn
}

func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

type recommend_gif struct {
	gifid string
	score float64
}

type recommend_gifs []recommend_gif

func (s recommend_gifs) Len() int           { return len(s) }
func (s recommend_gifs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s recommend_gifs) Less(i, j int) bool { return s[i].score > s[j].score }

func UserCF(likes map[string][]string, likes_u2g map[string][]string) map[string][]string {
	var mat map[string]map[string]float64
	mat = make(map[string]map[string]float64)
	var score float64
	for user1 := range likes_u2g {
		score = 0
		mat[user1] = make(map[string]float64)
		for user2 := range likes_u2g {
			intersect_like := intersect(likes_u2g[user1], likes_u2g[user2])
			for _, gifid := range intersect_like {
				score += 1 / (math.Log(float64(1 + len(likes[gifid]))))
			}
			score = score / math.Sqrt(float64(len(likes_u2g[user1])*len(likes_u2g[user2])))
			mat[user1][user2] = score
		}
	}

	var recommend_user map[string][]recommend_gif
	recommend_user = make(map[string][]recommend_gif)
	for user := range likes_u2g {
		recommend_user[user] = make([]recommend_gif, 0)
		for gifid := range likes {
			userCF_gif := recommend_gif{gifid: gifid, score: 0}
			for _, usr := range likes[gifid] {
				if usr != user {
					userCF_gif.score += mat[user][usr]
				}
			}
			recommend_user[user] = append(recommend_user[user], userCF_gif)
		}
	}

	// fmt.Println(recommend_user)
	for user := range likes_u2g {
		sort.Sort(recommend_gifs(recommend_user[user]))
	}

	var recommendation map[string][]string
	recommendation = make(map[string][]string)
	for user := range likes_u2g {
		recommendation[user] = make([]string, 0)
		for _, gif := range recommend_user[user] {
			recommendation[user] = append(recommendation[user], gif.gifid)
		}
		recommendation[user] = difference(recommendation[user], likes_u2g[user])
		if 10 < len(recommendation[user]) {
			recommendation[user] = recommendation[user][:10]
		}
	}

	return recommendation
}
