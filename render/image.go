package render

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"log"
	"strings"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const FontSize = 12

var dejaVuMonoTtf, dejaVuMonoBoldTtf *truetype.Font

func init() {
	var err error

	if dejaVuMonoTtf, err = freetype.ParseFont(dejaVuMono); err != nil {
		log.Panicf("Unable to parse DejaVu Mono: %v", err)
	}

	if dejaVuMonoBoldTtf, err = freetype.ParseFont(dejaVuMonoBold); err != nil {
		log.Panicf("Unable to parse DejaVu Mono Bold: %v", err)
	}
}

func RenderImage(tmpl string) (string, error) {
	// Set up font first so we know the height.
	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(dejaVuMonoTtf)
	ctx.SetFontSize(FontSize)
	lineHeight := ctx.PointToFix32(FontSize)
	ctx.SetHinting(freetype.FullHinting)
	// Create image.
	lines := strings.Split(tmpl, "\n")
	maxX := 400
	maxY := (len(lines) + 1) * int(lineHeight>>8)
	m := image.NewRGBA(image.Rect(0, 0, maxX, maxY))
	draw.Draw(m, m.Bounds(), image.White, image.ZP, draw.Src)
	// Bind context to image.
	ctx.SetClip(m.Bounds())
	ctx.SetDst(m)
	// Render lines.
	for i, l := range lines {
		ctx.SetSrc(image.Black)
		ctx.DrawString(l, raster.Point{0, lineHeight * raster.Fix32(i+1)})
	}
	// Output
	buf := bytes.NewBuffer([]byte{})
	err := png.Encode(buf, m)
	return buf.String(), err
}
