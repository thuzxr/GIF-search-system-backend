package recommend

import(
	"testing"

	"github.com/stretchr/testify/assert"
	"backend/utils"
)

func TestRecommend(t *testing.T){
	gifs:=utils.JsonParse("../info_old_recommend.json");
	rec_gif:=Recommend(gifs[1], gifs)
	assert.Equal(t, len(rec_gif), 10)
}