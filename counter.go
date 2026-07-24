package main

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const counterFontSize = 24

func NewCounter(jellyFishCount int) *Counter {
	// counterFaceSource holds the parsed font; parsing is expensive so it
	// happens once at init, never inside Draw.
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}

	return &Counter{
		unitsEaten:          0,
		totalJellyFishCount: jellyFishCount,
		fontSource:          s,
	}
}

type Counter struct {
	unitsEaten          int
	totalJellyFishCount int
	fontSource          *text.GoTextFaceSource
}

func (c *Counter) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf("%d/%d", c.unitsEaten, c.totalJellyFishCount)

	face := &text.GoTextFace{
		Source: c.fontSource,
		Size:   counterFontSize,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64((screenWidth/2)-len(msg)), 10)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, msg, face, op)
}

func (c *Counter) Add(val int) {
	c.unitsEaten += val
}
