package recommend

import(
	"testing"

	"github.com/stretchr/testify/assert"
	"backend/utils"
	"backend/database"
)

func TestRecommend(t *testing.T){
	gifs:=utils.JsonParse("../info_old_recommend.json");
	rec_gif:=Recommend(gifs[1], gifs)
	assert.Equal(t, len(rec_gif), 10)
}

func TestUserCF(t *testing.T){
	DB:=database.ConnectDB()
	_, _, _, likes, likes_u2g := database.LoadAll(DB)
	_ = UserCF(likes, likes_u2g)
}