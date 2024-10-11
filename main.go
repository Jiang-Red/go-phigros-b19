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

	"github.com/Jiang-Red/go-phigros-b19/phigros/phigros"

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
	err = Renderb19(j.PlayerInfo.Name, strconv.FormatFloat(float64(j.Summary.Rks), 'f', 6, 64), challengemoderank[(j.Summary.ChallengeModeRank-(j.Summary.ChallengeModeRank%100))/100], strconv.Itoa(int(j.Summary.ChallengeModeRank%100)), "10001", j.ScoreAcc)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(now).String())
}

var filepath = "D:/!!!important/go-phigros-b19/res/"

// Renderb19 ...
func Renderb19(plname, allrks, chal, chalnum, uid string, list []phigros.ScoreAcc) (err error) {
	errs := make(chan error)
	const w, h = 2360, 4780
	canvas := gg.NewContext(w, h)
	//canvas.SetRGB255(0, 255, 0)
	//canvas.Clear()
	a := 75.0

	x, y := 188.0, 682.0
	xspac, yspac := 1090.0, 160
	//var xj, yj float64 = 1090, 160

	cardimgs := make([]image.Image, len(list))

	wg := &sync.WaitGroup{}
	wg.Add(len(list) + 1)
	for i := 0; i < len(list); i++ {
		go func(i int) {
			defer wg.Done()
			cardimgs[i], err = drawcardback(canvas.W(), canvas.H(), i, a, x, y, list[i])
			if err != nil {
				errs <- err
				return
			}
		}(i)
	}

	drawfile, err := os.ReadDir(filepath + Illustration)
	if err != nil {
		errs <- err
		return
	}
	imgs, err := gg.LoadImage(filepath + Illustration + drawfile[rand.Intn(len(drawfile))].Name())
	if err != nil {
		errs <- err
		return
	}

	blured := imaging.Blur(imgs, 30)

	//更改图片原点
	canvas.ScaleAbout(float64(h)/float64(blured.Bounds().Dy()), float64(h)/float64(blured.Bounds().Dy()), float64(w)/2, 0)

	canvas.DrawImageAnchored(blured, w/2, 0, 0.5, 0)

	canvas.Identity()

	go func() {
		defer wg.Done()
		drawParallelogram(canvas, a, 0, 166, 1324, 410) // h = 396
		canvas.SetRGBA255(0, 0, 0, 160)
		canvas.Fill()

		drawParallelogram(canvas, a, 1318, 192, 1200, 350) // h = 338
		canvas.SetRGBA255(0, 0, 0, 160)
		canvas.Fill()

		drawParallelogram(canvas, a, 1320, 164, 6, 414)
		canvas.SetRGBA255(255, 255, 255, 255)
		canvas.Fill()

		tw, th := drawParallelogram(canvas, a, 534, 4342, 1312, 342)
		canvas.SetRGBA255(0, 0, 0, 160)
		canvas.Fill()

		drawParallelogram(canvas, a, 530, 4340, 6, 346)
		canvas.SetRGBA255(255, 255, 255, 255)
		canvas.Fill()

		drawParallelogram(canvas, a, 1842, 4340, 6, 346)
		canvas.SetRGBA255(255, 255, 255, 255)
		canvas.Fill()
		err = canvas.ParseFontFace(fontsd, 60)
		if err != nil {
			errs <- err
			return
		}

		canvas.DrawStringAnchored("Create By ZeroBot-Plugin", w/2-tw/2, 4342+th/4, 0.5, 0.5)
		canvas.DrawStringAnchored("UI Designer: eastown", w/2-tw/2, 4342+th*2/4, 0.5, 0.5)
		canvas.DrawStringAnchored("*Phigros B19 Picture*", w/2-tw/2, 4342+th*3/4, 0.5, 0.5)

		logo, err := gg.LoadImage(filepath + "/" + uid + "/avatar.png")
		if err != nil {
			errs <- err
			return
		}
		canvas.DrawRoundedRectangle(50, 166+396/2-290/2, 290, 290, 30)
		canvas.Clip()
		canvas.ScaleAbout(290.0/float64(logo.Bounds().Dx()), 290.0/float64(logo.Bounds().Dy()), 50, 166+396/2)
		canvas.DrawImageAnchored(logo, 50, 166+396/2, 0, 0.5)
		canvas.Identity()
		canvas.ResetClip()

		err = canvas.ParseFontFace(fontsd, 90)
		if err != nil {
			errs <- err
			return
		}
		canvas.DrawStringAnchored("Phigros", 50+290+50, 166+396/3, 0, 0.5)
		canvas.DrawStringAnchored("RankingScore查询", 50+290+50, 166+396*2/3, 0, 0.5)

		err = canvas.ParseFontFace(fontsd, 54)
		if err != nil {
			errs <- err
			return
		}

		canvas.DrawStringAnchored("Player: "+plname, w-920, 192+338/4, 0, 0.5)
		canvas.DrawStringAnchored("RankingScore: "+allrks, w-920, 192+338*2/4, 0, 0.5)
		canvas.DrawStringAnchored("ChallengeMode: ", w-920, 192+338*3/4, 0, 0.5)
		if chal != "" {
			chall, err := gg.LoadPNG(filepath + Challengemode + chal + ".png")
			if err != nil {
				errs <- err
				return
			}

			challengemodew, _ := canvas.MeasureString("ChallengeMode: ")
			canvas.ScaleAbout(208.0/float64(chall.Bounds().Dx()), 100.0/float64(chall.Bounds().Dy()), w-920+challengemodew, 192+338*3/4)
			canvas.DrawImageAnchored(chall, w-920+int(challengemodew), 192+338*3/4, 0, 0.5)
			canvas.Identity()
			canvas.DrawStringAnchored(chalnum, w-920+challengemodew+208/2, 192+338*3/4, 0.5, 0.5)
		}
	}()
	wg.Wait()
	wg.Add(len(list))
	for i := 0; i < len(cardimgs); i++ {
		go func(i int) {
			defer wg.Done()
			canvas.DrawImage(cardimgs[i], 0, 0)
		}(i)
	}
	wg.Wait()
	// 画排名
	err = canvas.ParseFontFace(fontsd, 30)
	if err != nil {
		errs <- err
		return
	}
	tw, th := cal(a, 44)

	canvas.SetRGBA255(0, 0, 0, 255)

	for i := 0; i < len(list); i++ {
		spac := float64(yspac * i)
		if (i+1)%2 == 0 {
			x += xspac
		}

		if i == 0 {
			canvas.DrawStringAnchored("Phi", x+70/2-tw/2, spac+y+th/2, 0.5, 0.5)
		} else {
			canvas.DrawStringAnchored("#"+strconv.Itoa(i), x+70/2-tw/2, spac+y+th/2, 0.5, 0.5)
		}
		x = 188
	}

	// 画分数
	err = canvas.ParseFontFace(fontsd, 50)
	if err != nil {
		errs <- err
		return
	}

	_, th = cal(a, 218)

	canvas.SetRGBA255(255, 255, 255, 255)

	for i := 0; i < len(list); i++ {
		spac := float64(yspac * i)
		if (i+1)%2 == 0 {
			x += xspac
		}

		scorestr := strconv.Itoa(list[i].Score)
		if len(scorestr) < 7 {
			for i := len(scorestr); i < 7; i++ {
				scorestr = "0" + scorestr
			}
		}
		if list[i].Score != 0 {
			canvas.DrawStringAnchored(scorestr, x+408+518/2, y+th/2+spac, 0.5, 0.5)
		} else {
			canvas.DrawStringAnchored("0000000", x+408+518/2, y+th/2+spac, 0.5, 0.5)
		}
		x = 188
	}

	// 画acc
	err = canvas.ParseFontFace(fontsd, 44)
	if err != nil {
		errs <- err
		return
	}
	for i := 0; i < len(list); i++ {
		spac := float64(yspac * i)
		if (i+1)%2 == 0 {
			x += xspac
		}

		canvas.SetRGBA255(255, 255, 255, 255)
		if list[i].Acc != 0 {
			canvas.DrawStringAnchored(strconv.FormatFloat(float64(list[i].Acc), 'f', 2, 64)+"%", x+408+518/2, y+th*7/8+spac, 0.5, 0.5)
		} else {
			canvas.DrawStringAnchored("00.00%", x+408+518/2, y+th*7/8+spac, 0.5, 0.5)
		}
		x = 188
	}

	// 画曲名
	err = canvas.ParseFontFace(fontsd, 32)
	if err != nil {
		errs <- err
		return
	}
	for i := 0; i < len(list); i++ {
		spac := float64(yspac * i)
		if (i+1)%2 == 0 {
			x += xspac
		}

		canvas.SetRGBA255(255, 255, 255, 255)
		if list[i].SongId != "" {
			canvas.DrawStringAnchored(strings.Split(list[i].SongId, ".")[0], x+408+518/2, y+th/4+spac, 0.5, 0.5)
		} else {
			canvas.DrawStringAnchored(" ", x+408+326/2, y+th/4+spac, 0.5, 0.5)
		}
		x = 188
	}

	// 画定数
	err = canvas.ParseFontFace(fontsd, 30)
	if err != nil {
		errs <- err
		return
	}
	tw, th = cal(a, 94)
	for i := 0; i < len(list); i++ {
		spac := float64(yspac * i)
		if (i+1)%2 == 0 {
			x += xspac
		}

		canvas.SetRGBA255(255, 255, 255, 255)
		if list[i].Level != "" {
			canvas.DrawStringAnchored(list[i].Level+" "+strconv.FormatFloat(float64(list[i].Difficulty), 'f', 1, 64), x-36-tw/2+138/2, y+139+th/4+spac, 0.5, 0.5)
		} else {
			canvas.DrawStringAnchored("SP ?", x-36-tw/2+138/2, y+139+th/4+spac, 0.5, 0.5)
		}
		x = 188
	}

	err = canvas.ParseFontFace(fontsd, 44)
	if err != nil {
		errs <- err
		return
	}
	for i := 0; i < len(list); i++ {
		spac := float64(yspac * i)
		if (i+1)%2 == 0 {
			x += xspac
		}

		canvas.SetRGBA255(255, 255, 255, 255)
		if list[i].Rks != 0 {
			canvas.DrawStringAnchored(strconv.FormatFloat(float64(list[i].Rks), 'f', 2, 64), x-36-tw/2+138/2, y+139+th*2/3+spac, 0.5, 0.5)
		} else {
			canvas.DrawStringAnchored("0.00", x-36-tw/2+138/2, y+139+th*3/4+spac, 0.5, 0.5)
		}
		x = 188
	}

	select {
	case err := <-errs:
		return err
	default:
		wg.Wait()
	}
	return canvas.SavePNG(filepath + uid + "/output.png")
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

func drawcardback(w, h, i int, a, x, y float64, list phigros.ScoreAcc) (img image.Image, err error) {
	y += float64(160 * i)
	if (i+1)%2 == 0 {
		x += 1090
	}
	canvas := gg.NewContext(w, h)
	// 画排名背景
	drawParallelogram(canvas, a, x, y, 70, 44) // h = 42
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	// 画分数背景
	drawParallelogram(canvas, a, x+408, y+12, 518, 218) // h = 210
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
	canvas.ScaleAbout(110.0/float64(rankim.Bounds().Dx()), 110.0/float64(rankim.Bounds().Dy()), x+412, y+88)
	canvas.DrawImageAnchored(rankim, int(x)+412, int(y)+88, 0, 0)
	canvas.Identity()

	// 画分数线
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.DrawRectangle(x+536, y+142, 326, 2)
	canvas.Fill()

	// 画图片
	drawParallelogram(canvas, a, x+68, y, 348, 238)
	canvas.Clip()
	canvas.SetRGBA255(0, 0, 255, 0)
	canvas.Fill()

	_, err = os.Stat(filepath + Illustration + list.SongId + ".0.png")
	if list.SongId != "" && (err == nil || os.IsExist(err)) {
		var imgs image.Image
		imgs, err = gg.LoadImage(filepath + Illustration + list.SongId + ".0.png")
		if err != nil {
			return
		}
		canvas.ScaleAbout(436/float64(imgs.Bounds().Dx()), 230/float64(imgs.Bounds().Dy()), x, y)
		canvas.DrawImageAnchored(imgs, int(x), int(y), 0, 0)
		canvas.Identity()
	}
	canvas.ResetClip()

	// 画定数背景
	drawParallelogram(canvas, a, x-36, y+139, 138, 94) // h = 90
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
	drawParallelogram(canvas, a, x+926, y+10, 6, 222)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()
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
