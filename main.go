package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lianhong2758/PhigrosAPI/phigros"
	"github.com/nfnt/resize"

	"github.com/FloatTech/gg"
	"github.com/disintegration/imaging"
)

/*
type result struct {
	ID       int64   `db:"id"`
	Songname string  `db:"songname"` // eg. Shadow
	Diff     string  `db:"diff"`     // AT
	Diffnum  float64 `db:"diffnum"`  // 15.9
	Score    int64   `db:"score"`    // 1000000
	Acc      float64 `db:"acc"`      // 100.00
	Rank     string  `db:"rank"`     // phi
	Rksm     float64 `db:"rksm"`     // 15.90
}
*/

var Session = "nkyjch88ydrg4js83bea9jyiw"

var challengemoderank = []string{"white", "green", "blue", "red", "gold", "rainbow"}

var fontsd, _ = os.ReadFile(filepath + Font)

var now time.Time

func main() {
	err := phigros.LoadDifficult("./difficulty.tsv")
	if err != nil {
		panic(err)
	}
	j := phigros.UserRecord{}
	data, _ := phigros.GetDataFormTap(phigros.UserMeUrl, Session) //获取id
	var um phigros.UserMe
	_ = json.Unmarshal(data, &um)
	//ums, _ := json.MarshalIndent(um, "  ", "  ")
	//fmt.Println(string(ums))

	j.PlayerInfo = &phigros.PlayerInfo{
		Name:      um.Nickname,
		CreatedAt: um.CreatedAt,
		UpdatedAt: um.UpdatedAt,
		Avatar:    um.Avatar,
	}
	data, _ = phigros.GetDataFormTap(phigros.SaveUrl, Session) //获取存档链接

	var gs phigros.GameSave
	_ = json.Unmarshal(data, &gs)
	f, _ := os.Create("./gamesave.json")
	f.Write(data)

	//fmt.Println(gs.Results[0].GameFile.URL)
	ScoreAcc, err := phigros.ParseStatsByUrl(gs.Results[0].GameFile.URL)
	if err != nil {
		panic(err)
	}
	j.ScoreAcc = phigros.BN(ScoreAcc, 21)
	j.Summary = phigros.ProcessSummary(gs.Results[0].Summary)
	//js, _ := json.MarshalIndent(j, "  ", "  ")
	//fmt.Println(string(js))
	_, err = os.Stat(filepath + "10001/avatar.png")
	if os.IsNotExist(err) {
		var response *http.Response

		response, err = http.Get(j.PlayerInfo.Avatar)

		if err != nil {
			panic(err)
		}
		data, err = io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		err = os.WriteFile(filepath+"10001/avatar.png", data, 0644)
		if err != nil {
			panic(err)
		}
	}

	now = time.Now()
	DrawB19(1, j, strconv.FormatFloat(float64(j.Summary.Rks), 'f', 6, 64), challengemoderank[(j.Summary.ChallengeModeRank-(j.Summary.ChallengeModeRank%100))/100], strconv.Itoa(int(j.Summary.ChallengeModeRank%100)), "10001")
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(now).String())
}

var filepath = "res/"

func DrawB19(accuracy float64, j phigros.UserRecord, allrks, chal, chalnum, uid string) (err error) {
	var (
		w, h = int(2360 * accuracy), int(4780 * accuracy)
		//斜角阿尔法
		a = 75.0
		//基准定位
		x, y = 178.0 * accuracy, 682.0 * accuracy
		//并发
		wg = &sync.WaitGroup{}
	)
	//背景
	canvas := gg.NewContext(w, h)
	drawfile, err := os.ReadDir(filepath + Illustration)
	if err != nil {
		return
	}
	imgs, err := gg.LoadImage(filepath + Illustration + drawfile[rand.Intn(len(drawfile))].Name())
	if err != nil {
		return
	}
	imgs = imaging.Blur(imgs, 30)
	//速度提升0.7s
	imgs = resize.Resize(0, uint(h), imgs, resize.Bilinear) //改比例
	canvas.DrawImage(imgs, -(imgs.Bounds().Dx()-w)/2, 0)
	//其余的平行四边形底色

	drawParallelogram(canvas, a, 0, 166*accuracy, 1324*accuracy, 410*accuracy) // h = 396
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	drawParallelogram(canvas, a, 1318*accuracy, 192*accuracy, 1200*accuracy, 350*accuracy) // h = 338
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()
	//白线条
	drawParallelogram(canvas, a, 1320*accuracy, 164*accuracy, 6*accuracy, 414*accuracy)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()
	//底边
	tw, th := drawParallelogram(canvas, a, 534*accuracy, 4342*accuracy, 1312*accuracy, 342*accuracy)
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	drawParallelogram(canvas, a, 530, 4340, 6, 346)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	drawParallelogram(canvas, a, 1842, 4340, 6, 346)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()
	//version
	_ = canvas.ParseFontFace(fontsd, 60*accuracy)
	canvas.DrawStringAnchored("Create By ZeroBot-Plugin", float64(w)/2-tw/2, 4342*accuracy+th/4, 0.5, 0.5)
	canvas.DrawStringAnchored("UI Designer: eastown", float64(w)/2-tw/2, 4342*accuracy+th*2/4, 0.5, 0.5)
	canvas.DrawStringAnchored("*Phigros B19 Picture*", float64(w)/2-tw/2, 4342*accuracy+th*3/4, 0.5, 0.5)

	//头图文字
	_ = canvas.ParseFontFace(fontsd, 90*accuracy)
	canvas.DrawStringAnchored("Phigros", (50+290+50)*accuracy, (166+396/3)*accuracy, 0, 0.5)
	canvas.DrawStringAnchored("RankingScore查询", (50 + 290 + 50), (166+396*2/3)*accuracy, 0, 0.5)

	_ = canvas.ParseFontFace(fontsd, 54*accuracy)
	canvas.DrawStringAnchored("Player: "+j.PlayerInfo.Name, float64(w)-920*accuracy, (192+338/4)*accuracy, 0, 0.5)
	canvas.DrawStringAnchored("RankingScore: "+allrks, float64(w)-920, (192+338*2/4)*accuracy, 0, 0.5)
	canvas.DrawStringAnchored("ChallengeMode: ", float64(w)-920, (192+338*3/4)*accuracy, 0, 0.5)
	if chal != "" {
		chall, err := gg.LoadPNG(filepath + Challengemode + chal + ".png")
		if err != nil {
			return err
		}
		challengemodew, _ := canvas.MeasureString("ChallengeMode: ")
		chall = resize.Resize(208, 100, chall, resize.Bilinear)
		canvas.DrawImageAnchored(chall, w+int(-920*accuracy+challengemodew), int((192+338*3/4)*accuracy), 0, 0.5)

		canvas.DrawStringAnchored(chalnum, float64(w)-920*accuracy+challengemodew+(208/2)*accuracy, (192+338*3/4)*accuracy, 0.5, 0.5)
	}
	wg.Add(len(j.ScoreAcc) + 1)
	//头图
	func() {
		defer wg.Done()
		logo, err := gg.LoadImage(filepath + uid + "/avatar.png")
		if err != nil {
			return
		}
		logo = resize.Resize(uint(290.0*accuracy), uint(290.0*accuracy), logo, resize.Bilinear) //改比例

		canvas.DrawRoundedRectangle(50*accuracy, (166+396/2-290/2)*accuracy, 290*accuracy, 290*accuracy, 30)
		canvas.Clip()
		canvas.DrawImage(logo, int(50*accuracy), int((166+396/2-290/2)*accuracy))
		canvas.ResetClip()
	}()
	for i := 0; i < len(j.ScoreAcc); i++ {
		go func(i int) {
			defer wg.Done()
			cardimg, err := drawcardback(accuracy, i, canvas.W()-int(x), a, j.ScoreAcc[i])
			if err != nil {
				return
			}
			canvas.DrawImage(cardimg, int(x+1090*float64(i%2)*accuracy), int(float64(160*i)*accuracy+y))
		}(i)
	}
	wg.Wait()
	return canvas.SavePNG(filepath + uid + "/output5.png")
}

// 绘制平行四边形 angle 角度 x, y 坐标 w 宽度 l 斜边长
func drawParallelogram(canvas *gg.Context, angle, x, y, w, l float64) (tw, th float64) {
	// 左上角为原点
	x0, y0 := x, y
	// 右上角
	x1, y1 := x+w, y
	// 右下角
	tw, th = l*(math.Cos(angle*math.Pi/180.0)), l*(math.Sin(angle*math.Pi/180.0))
	x2 := x1 - tw
	y2 := y1 + th
	// 左下角
	x3, y3 := x2-w, y2
	canvas.NewSubPath()
	canvas.MoveTo(x0, y0)
	canvas.LineTo(x1, y1)
	canvas.LineTo(x2, y2)
	canvas.LineTo(x3, y3)
	canvas.ClosePath()
	return
}

func cal(angle, l float64) (tw, th float64) {
	return l * (math.Cos(angle * math.Pi / 180.0)), l * (math.Sin(angle * math.Pi / 180.0))
}

// accuracy精度,i索引,w宽度
func drawcardback(accuracy float64, i, w int, a float64, list phigros.ScoreAcc) (img image.Image, err error) {
	//xspac, yspac := 1090.0*accuracy, 160*accuracy
	iw := float64(100)
	canvas := gg.NewContext(w/2, int(220*accuracy))
	// 画排名背景
	drawParallelogram(canvas, a, iw*accuracy, 0, 70*accuracy, 44*accuracy) // h = 42
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	// 画分数背景
	drawParallelogram(canvas, a, (iw+408)*accuracy, 12*accuracy, 518*accuracy, 218*accuracy) // h = 210
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	// 画rank图标
	rank := ""
	if list.Fc && list.Score != 1000000 {
		rank = "fc"
	} else {
		rank = checkrank(int64(list.Score))
	}
	rankim, err := gg.LoadImage(filepath + Rank + rank + ".png")
	if err != nil {
		return
	}
	rankim = resize.Resize(uint(110*accuracy), uint(110*accuracy), rankim, resize.Bilinear)
	canvas.DrawImage(rankim, int((iw+412)*accuracy), int(88*accuracy))

	// 画分数线
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.DrawRectangle((536+iw)*accuracy, 142*accuracy, 326*accuracy, 2*accuracy)
	canvas.Fill()

	// 画图片
	drawParallelogram(canvas, a, (68+iw)*accuracy, 0, 348, 238)
	canvas.Clip()

	_, err = os.Stat(filepath + Illustration + list.SongId + ".0.png")
	if list.SongId != "" && (err == nil || os.IsExist(err)) {
		var imgs image.Image
		imgs, err = gg.LoadImage(filepath + Illustration + list.SongId + ".0.png")
		if err != nil {
			return
		}
		imgs = resize.Resize(uint(436*accuracy), uint(230*accuracy), imgs, resize.Bilinear)
		canvas.DrawImage(imgs, int(iw), 0)
	}
	canvas.ResetClip()

	// 画定数背景
	drawParallelogram(canvas, a, 36*accuracy, 139*accuracy, 138*accuracy, 94*accuracy) // h = 90
	switch list.Level {
	case "AT":
		canvas.SetRGBA255(56, 56, 56, 255)
	case "IN":
		canvas.SetRGBA255(190, 45, 35, 255)
	case "HD":
		canvas.SetRGBA255(3, 115, 190, 255)
	case "EZ":
		canvas.SetRGBA255(15, 180, 145, 255)
	default:
		canvas.SetRGBA255(56, 56, 56, 255)
	}
	canvas.Fill()

	// 画边缘
	drawParallelogram(canvas, a, 926*accuracy, 10*accuracy, 6*accuracy, 222*accuracy)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()
	//画文字

	// 画排名

	_ = canvas.ParseFontFace(fontsd, 30*accuracy)
	tw, th := cal(a, 44)
	canvas.SetRGBA255(0, 0, 0, 255)

	if i == 0 {
		canvas.DrawStringAnchored("Phi", (iw+(70/2)-tw/2)*accuracy, 0+th/2, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("#"+strconv.Itoa(i), (iw+70/2-tw/2)*accuracy, 0+th/2, 0.5, 0.5)
	}

	// 画分数
	_ = canvas.ParseFontFace(fontsd, 50*accuracy)

	_, th = cal(a, 218)

	canvas.SetRGBA255(255, 255, 255, 255)
	scorestr := strconv.Itoa(list.Score)
	if len(scorestr) < 7 {
		for i := len(scorestr); i < 7; i++ {
			scorestr = "0" + scorestr
		}
	}
	if list.Score != 0 {
		canvas.DrawStringAnchored(scorestr, (iw+408+518/2)*accuracy, th/2*accuracy, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("0000000", (iw+408+518/2)*accuracy, th/2*accuracy, 0.5, 0.5)
	}

	// 画acc
	_ = canvas.ParseFontFace(fontsd, 44*accuracy)
	if list.Acc != 0 {
		canvas.DrawStringAnchored(strconv.FormatFloat(float64(list.Acc), 'f', 2, 64)+"%",  (iw+408+518/2)*accuracy, th*7/8*accuracy, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("00.00%",  (iw+408+518/2)*accuracy, th*7/8*accuracy, 0.5, 0.5)
	}

	// 画曲名
	_ = canvas.ParseFontFace(fontsd, 32*accuracy)
	if list.SongId != "" {
		canvas.DrawStringAnchored(strings.Split(list.SongId, ".")[0],  (iw+408+518/2)*accuracy, th/4*accuracy, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored(" ",  (iw+408+326/2)*accuracy, th/4*accuracy, 0.5, 0.5)
	}
	// 画定数
	_ = canvas.ParseFontFace(fontsd, 30*accuracy)
	tw, th = cal(a, 94)
	if list.Level != "" {
		canvas.DrawStringAnchored(list.Level+" "+strconv.FormatFloat(float64(list.Difficulty), 'f', 1, 64),  (iw-36-tw/2+50)*accuracy, (139+th/4)*accuracy, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("SP ?",  (iw-36-tw/2+50)*accuracy, (139+th/4)*accuracy, 0.5, 0.5)
	}

	_ = canvas.ParseFontFace(fontsd, 44*accuracy)
	if list.Rks != 0 {
		canvas.DrawStringAnchored(strconv.FormatFloat(float64(list.Rks), 'f', 2, 64),  (iw-36-tw/2+40)*accuracy, (139+th*2/3)*accuracy, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("0.00",  (iw-36-tw/2+40)*accuracy, (139+th*2/3)*accuracy, 0.5, 0.5)
	}
	img = canvas.Image()
	return
}

func checkrank(score int64) string {
	if score == 1000000 {
		return "phi"
	}
	if score >= 960000 {
		return "v"
	}
	if score >= 920000 {
		return "s"
	}
	if score >= 880000 {
		return "a"
	}
	if score >= 820000 {
		return "b"
	}
	if score >= 700000 {
		return "c"
	}
	return "f"
}
