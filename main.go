package main

import (
	"encoding/json"
	"image"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"test/phigros/phigros"

	"github.com/FloatTech/floatbox/file"
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

var Session = "-"

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

	j.PlayerInfo = phigros.PlayerInfo{
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
	//js, _ := json.MarshalIndent(j, "  ", "  ")
	//fmt.Println(string(js))

	tb := 0.0
	for _, s := range j.ScoreAcc[:20] {
		tb += float64(s.Rks)
	}

	err = renderb19(j.PlayerInfo.Name, strconv.FormatFloat(tb/20, 'f', 2, 64), "Gold", "45", "10001", j.ScoreAcc)
	if err != nil {
		panic(err)
	}
}

var filepath = "D:/!!!important/go-phigros-b19/res/"

func renderb19(plname, allrks, chal, chalnum, uid string, list []phigros.ScoreAcc) error {
	const w, h = 2360, 4780
	canvas := gg.NewContext(w, h)
	//canvas.SetRGB255(0, 255, 0)
	//canvas.Clear()

	drawfile, err := os.ReadDir(filepath + Illustration)
	if err != nil {
		return err
	}

	imgs, err := gg.LoadImage(filepath + Illustration + drawfile[rand.Intn(len(drawfile))].Name())
	if err != nil {
		return err
	}

	blured := imaging.Blur(imgs, 30)

	a := 75.0

	canvas.ScaleAbout(9064.0/float64(blured.Bounds().Dx()), float64(h)/float64(blured.Bounds().Dy()), float64(w)/2, 0)

	canvas.DrawImageAnchored(blured, w/2, 0, 0.5, 0)

	canvas.Identity()

	drawParallelogram(canvas, a, 0, 166, 1324, 410) // h = 396
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	drawParallelogram(canvas, a, 1318, 192, 1200, 350) // h = 338
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	drawParallelogram(canvas, a, 1320, 164, 6, 414)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	drawParallelogram(canvas, a, 534, 4342, 1312, 342)
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	drawParallelogram(canvas, a, 530, 4340, 6, 346)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	drawParallelogram(canvas, a, 1842, 4340, 6, 346)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	font, err := gg.LoadFontFace(filepath+Font, 60)
	if err != nil {
		return err
	}
	canvas.SetFontFace(font)

	canvas.DrawStringAnchored("Create By ZeroBot-Plugin", w/2, 4342+346/4, 0.5, 0.5)
	canvas.DrawStringAnchored("UI Designer: eastown", w/2, 4342+346*2/4, 0.5, 0.5)
	canvas.DrawStringAnchored("*Phigros B19 Picture*", w/2, 4342+346*3/4, 0.5, 0.5)

	logo, err := gg.LoadImage(filepath + Icon)
	if err != nil {
		return err
	}
	canvas.ScaleAbout(290.0/float64(logo.Bounds().Dx()), 290.0/float64(logo.Bounds().Dy()), 50, 166+396/2)
	canvas.DrawImageAnchored(logo, 50, 166+396/2, 0, 0.5)
	canvas.Identity()

	font, err = gg.LoadFontFace(filepath+Font, 90)
	if err != nil {
		return err
	}
	canvas.SetFontFace(font)
	canvas.DrawStringAnchored("Phigros", 50+290+50, 166+396/3, 0, 0.5)
	canvas.DrawStringAnchored("RankingScore查询", 50+290+50, 166+396*2/3, 0, 0.5)

	font, err = gg.LoadFontFace(filepath+Font, 54)
	if err != nil {
		return err
	}
	canvas.SetFontFace(font)
	canvas.DrawStringAnchored("Player: "+plname, w-920, 192+338/4, 0, 0.5)
	canvas.DrawStringAnchored("RankingScore: "+allrks, w-920, 192+338*2/4, 0, 0.5)
	canvas.DrawStringAnchored("ChallengeMode: ", w-920, 192+338*3/4, 0, 0.5)
	if chal != "" {
		chall, err := gg.LoadPNG(filepath + Challengemode + chal + ".png")
		if err != nil {
			return err
		}

		challengemodew, _ := canvas.MeasureString("ChallengeMode: ")
		canvas.ScaleAbout(208.0/float64(chall.Bounds().Dx()), 100.0/float64(chall.Bounds().Dy()), w-920+challengemodew, 192+338*3/4)
		canvas.DrawImageAnchored(chall, w-920+int(challengemodew), 192+338*3/4, 0, 0.5)
		canvas.Identity()
		canvas.DrawStringAnchored(chalnum, w-920+challengemodew+208/2, 192+338*3/4, 0.5, 0.5)
	}

	var x, y float64 = 188, 682
	var i int64
	var xj, yj float64 = 1090, 160

	err = mix(canvas, i, a, x, y, list[i])
	if err != nil {
		return err
	}
	i++
	x += xj
	y += yj
	for ; i < 22; i++ {
		if i%2 == 0 {
			err := mix(canvas, i, a, x, y, list[i])
			if err != nil {
				return err
			}

			x += xj
			y += yj
		} else {
			err := mix(canvas, i, a, x, y, list[i])
			if err != nil {
				return err
			}

			x -= xj
			y += yj
		}
	}
	_ = os.Mkdir(filepath+uid, 0644)
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

func mix(canvas *gg.Context, i int64, a, x, y float64, list phigros.ScoreAcc) (err error) {
	// 画排名背景
	tw, th := drawParallelogram(canvas, a, x, y, 70, 44) // h = 42
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	// 画排名
	font, err := gg.LoadFontFace(filepath+Font, 30)
	if err != nil {
		return
	}
	canvas.SetFontFace(font)
	canvas.SetRGBA255(0, 0, 0, 255)

	if i == 0 {
		canvas.DrawStringAnchored("Phi", x-tw/2+70/2, y+th/2, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("#"+strconv.FormatInt(i, 10), x-tw/2+70/2, y+th/2, 0.5, 0.5)

	}

	// 画分数背景
	_, th = drawParallelogram(canvas, a, x+408, y+12, 518, 218) // h = 210
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

	// 画分数
	font, err = gg.LoadFontFace(filepath+Font, 50)
	if err != nil {
		return
	}
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	scorestr := strconv.Itoa(list.Score)
	if len(scorestr) < 7 {
		for i := len(scorestr); i < 7; i++ {
			scorestr = "0" + scorestr
		}
	}
	if list.Score != 0 {
		canvas.DrawStringAnchored(scorestr, x+408+518/2, y+th/2, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("0000000", x+408+518/2, y+th/2, 0.5, 0.5)
	}

	// 画acc
	font, err = gg.LoadFontFace(filepath+Font, 44)
	if err != nil {
		return
	}
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	if list.Acc != 0 {
		canvas.DrawStringAnchored(strconv.FormatFloat(float64(list.Acc), 'f', 2, 64)+"%", x+408+518/2, y+th*7/8, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("00.00%", x+408+518/2, y+th*7/8, 0.5, 0.5)
	}

	// 画曲名
	font, err = gg.LoadFontFace(filepath+Font, 32)
	if err != nil {
		return
	}
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	if list.SongId != "" {
		canvas.DrawStringAnchored(strings.Split(list.SongId, ".")[0], x+408+518/2, y+th/4, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored(" ", x+408+326/2, y+th/4, 0.5, 0.5)
	}

	// 画图片
	drawParallelogram(canvas, a, x+68, y, 348, 238)
	canvas.Clip()
	canvas.SetRGBA255(0, 0, 255, 0)
	canvas.Fill()

	if list.SongId != "" && file.IsExist(filepath+Illustration+strings.Split(list.SongId, ".")[0]+".png") {
		var imgs image.Image
		imgs, err = gg.LoadImage(filepath + Illustration + strings.Split(list.SongId, ".")[0] + ".png")
		if err != nil {
			return
		}
		canvas.ScaleAbout(436/float64(imgs.Bounds().Dx()), 230/float64(imgs.Bounds().Dy()), x, y)
		canvas.DrawImageAnchored(imgs, int(x), int(y), 0, 0)
		canvas.Identity()
	}
	canvas.ResetClip()

	// 画定数背景
	tw, th = drawParallelogram(canvas, a, x-36, y+139, 138, 94) // h = 90
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

	// 画定数
	font, err = gg.LoadFontFace(filepath+Font, 30)
	if err != nil {
		return
	}
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	if list.Level != "" {
		canvas.DrawStringAnchored(list.Level+" "+strconv.FormatFloat(float64(list.Difficulty), 'f', 1, 64), x-36-tw/2+138/2, y+139+th/4, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("SP ?", x-36-tw/2+138/2, y+139+th/4, 0.5, 0.5)
	}

	font, err = gg.LoadFontFace(filepath+Font, 44)
	if err != nil {
		return
	}
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	if list.Rks != 0 {
		canvas.DrawStringAnchored(strconv.FormatFloat(float64(list.Rks), 'f', 2, 64), x-36-tw/2+138/2, y+139+th*2/3, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored("0.00", x-36-tw/2+138/2, y+139+th*3/4, 0.5, 0.5)
	}

	// 画边缘
	drawParallelogram(canvas, a, x+926, y+10, 6, 222)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()
	return nil
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
