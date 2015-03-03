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

const (
	FontSize = 13
	Spacing  = 4
)

var dejaVuMonoTtf, dejaVuMonoBoldTtf *truetype.Font

var imageColours = map[string]image.Image{
	Black:   image.NewUniform(HtmlRgbColours[Black].ToRgba()),
	Red:     image.NewUniform(HtmlRgbColours[Red].ToRgba()),
	Green:   image.NewUniform(HtmlRgbColours[Green].ToRgba()),
	Yellow:  image.NewUniform(HtmlRgbColours[Yellow].ToRgba()),
	Blue:    image.NewUniform(HtmlRgbColours[Blue].ToRgba()),
	Magenta: image.NewUniform(HtmlRgbColours[Magenta].ToRgba()),
	Cyan:    image.NewUniform(HtmlRgbColours[Cyan].ToRgba()),
	Gray:    image.NewUniform(HtmlRgbColours[Gray].ToRgba()),
}

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
	lineHeight := ctx.PointToFix32(FontSize + Spacing)
	// Create image.
	lines := strings.Split(tmpl, "\n")
	maxX := 400
	maxY := (len(lines) + 1) * int(lineHeight>>8)
	m := image.NewRGBA(image.Rect(0, 0, maxX, maxY))
	draw.Draw(m, m.Bounds(), image.White, image.ZP, draw.Src)
	// Bind context to image.
	ctx.SetClip(m.Bounds())
	ctx.SetDst(m)
	// Track position and line
	x := raster.Fix32(0)
	line := raster.Fix32(1)
	WalkTemplate(tmpl, func(text, colour string, bold bool) {
		ctx.SetSrc(imageColours[colour])
		if bold {
			ctx.SetFont(dejaVuMonoBoldTtf)
		} else {
			ctx.SetFont(dejaVuMonoTtf)
		}
		lines := strings.Split(text, "\n")
		numLines := len(lines)
		for i, l := range lines {
			if len(l) > 0 {
				point, err := ctx.DrawString(l, raster.Point{x, lineHeight * line})
				if err != nil {
					log.Panicf("Error writing: %v\nText: %#v", err, l)
				}
				x = point.X
			}
			if i < numLines-1 {
				line++
				x = 0
			}
		}
	})
	// Output
	buf := bytes.NewBuffer([]byte{})
	err := png.Encode(buf, m)
	return buf.String(), err
}
