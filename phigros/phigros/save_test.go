package phigros_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/lianhong2758/PhigrosAPI/phigros"
)

// 个人Session获取查看link: https://www.taptap.cn/moment/535045245566452043
// 2.2
var Session = "zp4fenkahnb7639x9t6v38jmn"

func TestSave(t *testing.T) {
	//data, _ := phigros.GetDataFormTap(phigros.UserMeUrl, Session) //获取id
	data, _ := phigros.GetDataFormTap(phigros.SaveUrl, Session) //获取存档链接
	var us phigros.GameSave
	_ = json.Unmarshal(data, &us)
	_ = os.MkdirAll("../data/gamesave/", os.ModePerm)
	phigros.SaveGameData(us.Results[0].GameFile.URL, "../data/gamesave/"+Session+".zip")
	_ = phigros.LoadDifficult("../difficulty.tsv")
	j, _ := phigros.ParseSave("../data/gamesave/" + Session + ".zip")
	fmt.Println(j)
}
func TestJson(t *testing.T) {
	_ = phigros.LoadDifficult("../difficulty.tsv")
	j := phigros.UserRecord{Session: Session}
	data, _ := phigros.GetDataFormTap(phigros.UserMeUrl, Session) //获取id
	var um phigros.UserMe
	_ = json.Unmarshal(data, &um)
	j.PlayerInfo = phigros.PlayerInfo{
		Name:      um.Nickname,
		CreatedAt: um.CreatedAt,
		UpdatedAt: um.UpdatedAt,
		Avatar:    um.Avatar,
	}
	data, _ = phigros.GetDataFormTap(phigros.SaveUrl, Session) //获取存档链接
	var gs phigros.GameSave
	_ = json.Unmarshal(data, &gs)
	ScoreAcc, _ := phigros.ParseStatsByUrl(gs.Results[0].GameFile.URL)
	j.ScoreAcc = phigros.BN(ScoreAcc, 10)
	fmt.Println(j)
}
