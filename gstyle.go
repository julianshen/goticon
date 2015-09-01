package main

import (
	"crypto/sha512"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"sort"
)

var (
	BACKGROUND_COLOR color.NRGBA = color.NRGBA{224, 224, 224, 255}
)

type SortedBytes []byte

func (s SortedBytes) Len() int {
	return len(s)
}

func (s SortedBytes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortedBytes) Less(i, j int) bool { return s[i] < s[j] }

func findMedian(data []byte) byte {
	d := make([]byte, len(data), len(data))
	copy(d, data)
	sort.Sort(SortedBytes(d))
	return d[len(d)/2]
}

func GenerateIdenticon(data []byte, width int, margin int) image.Image {
	sum := sha512.Sum512(data)
	median := findMedian(sum[0:25])

	rand.Seed(int64(sum[26]))
	fg := colorful.FastWarmColor()
	palette := color.Palette{BACKGROUND_COLOR, fg}
	img := image.NewPaletted(image.Rect(0, 0, width, width), palette)

	u := image.NewUniform(fg)

	cubeWidth := (width - margin*2) / 5

	for i, v := range sum {
		if v >= median && i < 25 {
			x := (i%5)*cubeWidth + margin
			y := (i/5)*cubeWidth + margin
			draw.Draw(img, image.Rect(x, y, x+cubeWidth, y+cubeWidth), u, image.Point{0, 0}, draw.Over)
		}
	}

	return img
}
