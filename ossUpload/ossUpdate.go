package ossUpload

import (
	"backend/utils"
)

func OssUpdate(gifs []utils.Gifs){
	for i:=range(gifs)	{
		gifs[i].Oss_url=OssSignLink(gifs[i], 3600);
	}
}