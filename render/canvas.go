package render

import (
	"bytes"
	"regexp"
	"strings"
)

var pixelRegexp = regexp.MustCompile(`^((?:\{\{.+?\}\})*)(.)((?:\{\{_.+?\}\})*)`)

// Canvas is a two dimensional plane to which you can write strings by location.
type Canvas struct {
	pixels map[int]map[int]string
}

// NewCanvas initialises a Canvas struct.
func NewCanvas() *Canvas {
	return &Canvas{
		pixels: map[int]map[int]string{},
	}
}

func readPixel(input string) (pixel, prefix, suffix, remaining string) {
	if input == "" {
		return "", "", "", ""
	}
	matches := pixelRegexp.FindStringSubmatch(input)
	if matches == nil {
		return "", "", "", ""
	}
	return matches[2], matches[1], matches[3], input[len(matches[0]):]
}

// Draw draws a string to the canvas at the specified location.
func (c *Canvas) Draw(x, y int, content string) {
	xPtr := x
	yPtr := y
	remaining := content
	for len(remaining) > 0 {
		var pixel, prefix, suffix string
		pixel, prefix, suffix, remaining = readPixel(remaining)
		if c.pixels[yPtr] == nil {
			c.pixels[yPtr] = map[int]string{}
		}
		c.pixels[yPtr][xPtr] = prefix + pixel + suffix
		xPtr += 1
	}
}

// Render render all of the drawn objects in the smallest possible rectangle.
// Render calculates the smallest rectangle containing all drawn objects and
// defers to RenderRect.
func (c *Canvas) Render() string {
	var minX, minY, maxX, maxY int
	firstPass := true
	for y, row := range c.pixels {
		if firstPass || y < minY {
			minY = y
		}
		if firstPass || y > maxY {
			maxY = y
		}
		for x, _ := range row {
			if firstPass || x < minX {
				minX = x
			}
			if firstPass || x > maxX {
				maxX = x
			}
			firstPass = false
		}
		firstPass = false
	}
	return c.RenderRect(minX, minY, maxX-minX+1, maxY-minY+1)
}

// RenderRect renders all of the drawn objects given the passed rectangle.
func (c *Canvas) RenderRect(x, y, width, height int) string {
	lines := []string{}
	for yPtr := y; yPtr < y+height; yPtr++ {
		if c.pixels[yPtr] == nil {
			lines = append(lines, strings.Repeat(" ", width))
			continue
		}
		line := bytes.Buffer{}
		for xPtr := x; xPtr < x+width; xPtr++ {
			pixel := c.pixels[yPtr][xPtr]
			if pixel == "" {
				pixel = " "
			}
			line.WriteString(pixel)
		}
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n")
}
