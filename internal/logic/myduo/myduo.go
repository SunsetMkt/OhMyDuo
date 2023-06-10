package myduo

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"my-duo/internal/consts"
	"my-duo/internal/service"
	"my-duo/internal/utils"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type sMyDuo struct {
	ImgSize         image.Rectangle
	FontUnicode     *truetype.Font
	FontAsciiBold   *truetype.Font
	BackgroundImage image.Image
	RoundMaskImage  image.Image
}

func prepareFont(path string) *truetype.Font {
	tt, err := freetype.ParseFont(utils.GetResource("resource/public/resource/font/" + path))
	if err != nil {
		panic(err)
	}
	return tt
}

func prepareImage(path string) image.Image {
	im, _, err := image.Decode(
		bytes.NewReader(utils.GetResource(
			"resource/public/resource/image/" + path)),
	)
	if err != nil {
		panic(err)
	}
	return im
}

func init() {
	service.RegisterMyDuo(New())
}
func New() *sMyDuo {
	return &sMyDuo{
		ImgSize:         image.Rect(0, 0, 793, 793),
		FontUnicode:     prepareFont("OPPOSans-R.ttf"),
		FontAsciiBold:   prepareFont("DIN Next Rounded LT W05 Bold.ttf"),
		BackgroundImage: prepareImage("background-image.png"),
		RoundMaskImage:  prepareImage("round-mask.png"),
	}
}

func (sv *sMyDuo) Draw(ctx context.Context, elem consts.MyDuoElements) []byte {
	img := image.NewRGBA(sv.ImgSize)
	sv.drawBackground(img)
	sv.drawCharacter(img, elem.Character)

	// posX := 80
	// posY := 290
	// bg := image.NewUniform(color.RGBA{75, 75, 75, 255})
	// point := fixed.Point26_6{X: fixed.Int26_6(posX * 64), Y: fixed.Int26_6(posY * 64)}
	// drawDst := image.NewRGBA(img.Bounds())
	// drawer := &font.Drawer{
	// 	Dst: drawDst,
	// 	Src: bg,
	// 	Face: truetype.NewFace(
	// 		sv.FontAsciiBold,
	// 		&truetype.Options{Size: 45}),
	// 	Dot: point,
	// }
	// drawer.DrawString(elem.OriginText)
	textImg := sv.drawText("2345678 345672asdew8 自慰")
	draw.Draw(img, img.Bounds(), textImg, textImg.Bounds().Min, draw.Over)

	// encode to bytes
	buff := new(bytes.Buffer)
	png.Encode(buff, img)
	return buff.Bytes()
}

func (sv *sMyDuo) drawBackground(img *image.RGBA) {
	draw.Draw(img, img.Bounds(), sv.BackgroundImage, sv.BackgroundImage.Bounds().Min, draw.Src)
}

func (sv *sMyDuo) drawCharacter(img *image.RGBA, character consts.MyDuoCharacters) {
	c := prepareImage("characters/" + character.ToString() + ".png")
	draw.Draw(img, img.Bounds(), c, c.Bounds().Min, draw.Over)
}

func (sv *sMyDuo) drawTextOnImg(img *image.RGBA, f font.Face, s string, x fixed.Int26_6, y fixed.Int26_6) {
	// max width = 635 px
	// text size = 64 px
	// line spacing = 18 px
	bg := image.NewUniform(color.RGBA{75, 75, 75, 255})
	xp, yp := int(x/64), int(y/64)
	println(xp)
	println(yp)
	point := fixed.Point26_6{X: 0, Y: 0}
	drawDst := image.NewRGBA(img.Bounds())
	drawer := &font.Drawer{
		Dst:  drawDst,
		Src:  bg,
		Face: f,
		Dot:  point,
	}
	drawer.DrawString(s)
	draw.Draw(img, img.Bounds(), drawDst, drawDst.Bounds().Min.Add(image.Point{xp, yp}), draw.Over)
}

func (sv *sMyDuo) drawText(s string) *image.RGBA {
	// max width = 635 px
	// text size = 64 px
	// line spacing = 18 px
	img := image.NewRGBA(image.Rect(0, 0, 635, 64+18+64))
	faceAsciiBold := truetype.NewFace(sv.FontAsciiBold, &truetype.Options{Size: 45})
	faceUnicode := truetype.NewFace(sv.FontUnicode, &truetype.Options{Size: 45})
	pieces := utils.SplitText(s)
	lines := 1
	cursor := fixed.Int26_6(0 * 64)
	max_width := fixed.Int26_6(635 * 64)
	for _, v := range pieces {
		var w fixed.Int26_6 = 0
		var f font.Face
		if v.Unicode {
			f = faceUnicode
		} else {
			f = faceAsciiBold
		}
		w = font.MeasureString(f, v.Text)
		// todo: split long words
		if cursor+w > max_width {
			// next line
			lines++
			cursor = 0
		}
		sv.drawTextOnImg(img, f, v.Text, cursor, fixed.Int26_6(lines*18*64))
		cursor += w
	}
	return img
	// font.MeasureString()
}
