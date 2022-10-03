package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/zbputils/img"
)

var (
	// 排名背景
	x, y float64 = 188, 682
	w, h float64 = 70, 44
	// 图片
	x1, y1 float64 = 256, 682
	w1, h1 float64 = 346, 238
	// 定数背景
	x2, y2 float64 = 152, 821
	w2, h2 float64 = 138, 94
	// 分数背景
	x3, y3 float64 = 596, 694
	w3, h3 float64 = 518, 218
	// 边缘
	x4, y4 float64 = 1114, 692
	w4, h4 float64 = 6, 222
	// 真图片
	x5, y5 int = 194, 682
	// 排名
	x6, y6 float64 = 178, 714
	// 分数线
	x7, y7 float64 = 724, 824
	w7, h7 float64 = 326, 2
)
var (
	diff  float64 = 15.90
	level string  = "AT 15.9"
	//level
	x8, y8 float64 = 144, 856
	level2 float64 = 15.90
	//level2
	x9, y9 float64 = 138, 898
	// rank
	x10, y10 int    = 600, 770
	score    string = "1000000"
	// score
	x11, y11 float64 = 776, 798
	name     string  = "Shadow"
	// name
	x12, y12 float64 = 594, 740
	acc      float64 = 100.00
	// acc
	x13, y13 float64 = 576, 878

	pl, rks, cm string = "Player: ", "RankingScore: ", "ChallengeMode: "
)

// 角度
var a float64 = 75

func main() {
	file, _ := getfile()
	fmt.Println(file)

	canvas := gg.NewContext(2360, 4780)
	canvas.SetRGB255(0, 255, 0)
	canvas.Clear()

	imgs, _ := img.LoadFirstFrame(file+Res+"/back.png", 2048, 1080)

	blured := imgs.Blur(30)

	cutted := cut4img(imgs, a)

	canvas.DrawImage(img.Size(blured.Im, 9064, 4780).Im, -3352, 0)

	draw4(canvas, a, 0, 166, 1324, 410)
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	draw4(canvas, a, 1318, 192, 1200, 350)
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	draw4(canvas, a, 1320, 164, 6, 414)
	canvas.SetColor(color.White)
	canvas.Fill()

	logo, _ := gg.LoadPNG(file + Icon)
	canvas.DrawImage(img.Size(logo, 290, 290).Im, 50, 216)

	font, _ := gg.LoadFontFace(file+Font, 90)
	canvas.SetFontFace(font)
	canvas.DrawString("Phigros", 422, 336)
	canvas.DrawString("RankingScore查询", 422, 462)

	font, _ = gg.LoadFontFace(file+Font, 54)
	canvas.SetFontFace(font)
	canvas.DrawString(pl+"yyw", 1434, 300)
	canvas.DrawString(rks+strconv.FormatFloat(allrks(rksc(acc, diff), rksc(acc, diff)), 'f', 3, 64), 1434, 380)

	chall, _ := gg.LoadPNG(file + Challengemode + "/rainbow.png")
	canvas.DrawImage(img.Size(chall, 208, 100).Im, 1848, 392)

	canvas.DrawString(cm, 1434, 460)
	canvas.DrawString("49", 1916, 460)

	rank, _ := img.LoadFirstFrame(file+Rank+"/phi.png", 110, 110)

	var i int64 = 0
	var xj, yj float64 = 1090, 200
	for ; i < 20; i++ {
		if i%2 == 0 {
			mix(canvas, i, file, cutted, rank)

			x, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11, x12, x13 = x+xj, x1+xj, x2+xj, x3+xj, x4+xj, x5+int(xj), x6+xj, x7+xj, x8+xj, x9+xj, x10+int(xj), x11+xj, x12+xj, x13+xj
			y, y1, y2, y3, y4, y5, y6, y7, y8, y9, y10, y11, y12, y13 = y+yj, y1+yj, y2+yj, y3+yj, y4+yj, y5+int(yj), y6+yj, y7+yj, y8+yj, y9+yj, y10+int(yj), y11+yj, y12+yj, y13+yj
		} else {
			mix(canvas, i, file, cutted, rank)

			x, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11, x12, x13 = x-xj, x1-xj, x2-xj, x3-xj, x4-xj, x5-int(xj), x6-xj, x7-xj, x8-xj, x9-xj, x10-int(xj), x11-xj, x12-xj, x13-xj
			y, y1, y2, y3, y4, y5, y6, y7, y8, y9, y10, y11, y12, y13 = y+yj, y1+yj, y2+yj, y3+yj, y4+yj, y5+int(yj), y6+yj, y7+yj, y8+yj, y9+yj, y10+int(yj), y11+yj, y12+yj, y13+yj
		}
	}
	x, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11, x12, x13 = x-xj, x1-xj, x2-xj, x3-xj, x4-xj, x5-int(xj), x6-xj, x7-xj, x8-xj, x9-xj, x10-int(xj), x11-xj, x12-xj, x13-xj
	y, y1, y2, y3, y4, y5, y6, y7, y8, y9, y10, y11, y12, y13 = y+-(yj*19), y1+-(yj*19), y2+-(yj*19), y3+-(yj*19), y4+-(yj*19), y5-int(yj*19), y6+-(yj*19), y7+-(yj*19), y8+-(yj*19), y9+-(yj*19), y10-int(yj*19), y11+-(yj*19), y12+-(yj*19), y13-(yj*19)

	canvas.SavePNG(file + "/output.png")
}

func getfile() (file string, err error) {
	file, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	file = strings.Replace(file, "\\", "/", -1)
	return
}

// 绘制平行四边形 angle 角度 x, y 坐标 w 宽度 l 斜边长
func draw4(canvas *gg.Context, angle, x, y, w, l float64) {
	// 左上角为原点
	x0, y0 := x, y
	// 右上角
	x1, y1 := x+w, y
	// 右下角
	x2 := x1 - (l * (math.Cos(angle * math.Pi / 180.0)))
	y2 := y1 + (l * (math.Sin(angle * math.Pi / 180.0)))
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

func mix(canvas *gg.Context, i int64, file string, imgs, rank *img.Factory) {
	// 画排名背景
	draw4(canvas, a, x, y, w, h)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()

	// 画排名
	font, _ := gg.LoadFontFace(file+Font, 30)
	canvas.SetFontFace(font)
	canvas.SetRGBA255(0, 0, 0, 255)
	var fw2 float64
	if i == 0 {
		fw2, _ = canvas.MeasureString("Phi")
		canvas.DrawString("Phi", x6+((w-fw2)/2), y6)
	} else {
		fw2, _ = canvas.MeasureString("#" + strconv.FormatInt(i, 10))
		canvas.DrawString("#"+strconv.FormatInt(i, 10), x6+((w-fw2)/2), y6)
	}

	// 画分数背景
	draw4(canvas, a, x3, y3, w3, h3)
	canvas.SetRGBA255(0, 0, 0, 160)
	canvas.Fill()

	// 画rank图标
	canvas.DrawImage(rank.Im, x10, y10)

	// 画分数线
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.DrawRectangle(x7, y7, w7, h7)
	canvas.Fill()

	// 画分数
	font, _ = gg.LoadFontFace(file+Font, 50)
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.DrawString(score, x11, y11)

	// 画acc
	font, _ = gg.LoadFontFace(file+Font, 44)
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	fw, _ := canvas.MeasureString(strconv.FormatFloat(acc, 'f', 2, 64) + "%")
	canvas.DrawString(strconv.FormatFloat(acc, 'f', 2, 64)+"%", x13+((w3-fw)/2), y13)

	// 画曲名
	font, _ = gg.LoadFontFace(file+Font, 32)
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	fw1, _ := canvas.MeasureString(name)
	canvas.DrawString(name, x12+((w3-fw1)/2), y12)

	// 画图片
	draw4(canvas, a, x1, y1, w1, h1)
	canvas.SetRGBA255(0, 0, 255, 0)
	canvas.DrawImage(img.Size(imgs.Im, 436, 230).Im, x5, y5)
	canvas.Fill()

	// 画定数背景
	draw4(canvas, a, x2, y2, w2, h2)
	canvas.SetRGBA255(255, 0, 0, 255)
	canvas.Fill()

	// 画定数
	font, _ = gg.LoadFontFace(file+Font, 30)
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	fw3, _ := canvas.MeasureString(level)
	canvas.DrawString(level, x8+((w2-fw3)/2), y8)

	font, _ = gg.LoadFontFace(file+Font, 44)
	canvas.SetFontFace(font)
	canvas.SetRGBA255(255, 255, 255, 255)
	fw4, _ := canvas.MeasureString(strconv.FormatFloat(rksc(acc, diff), 'f', 2, 64))
	canvas.DrawString(strconv.FormatFloat(rksc(acc, diff), 'f', 2, 64), x9+((w2-fw4)/2), y9)

	// 画边缘
	draw4(canvas, a, x4, y4, w4, h4)
	canvas.SetRGBA255(255, 255, 255, 255)
	canvas.Fill()
}

// 将矩形图片裁切为平行四边形 angle 为角度
func cut4img(imgs *img.Factory, angle float64) *img.Factory {
	db := imgs.Im.Bounds()
	dst := imgs
	maxy := db.Max.Y
	maxx := db.Max.X
	sax := (float64(maxy) * (math.Cos(angle * math.Pi / 180.0)))
	ax := sax
	for autoadd := 1; autoadd < maxy; autoadd++ {
		for ; ax > 0; ax-- {
			dst.Im.Set(int(ax), int(autoadd), color.NRGBA{0, 0, 0, 0})
			dst.Im.Set(maxx+int(-ax), maxy-autoadd, color.NRGBA{0, 0, 0, 0})
		}
		ax = (float64(maxy-autoadd) * (math.Cos(angle * math.Pi / 180.0)))
	}
	return dst
}

func rksc(accc, diff float64) float64 {
	return ((100.0*(acc/100.0) - 55.0) / 45.0) * ((100.0*(acc/100.0) - 55.0) / 45.0) * diff
}

func allrks(phi, sco float64) float64 {
	return (phi + (sco * 19)) / 20
}
